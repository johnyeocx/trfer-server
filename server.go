package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/db"
	my_fb "github.com/johnyeocx/usual/server2/firebase"
	"github.com/johnyeocx/usual/server2/routes"
	"github.com/johnyeocx/usual/server2/utils/cloud"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env file
	err := godotenv.Load(".env")
	if err != nil {
		return 
	}

	// 2. Connect to services
	psqlDB := db.Connect()
	s3Cli := cloud.ConnectAWS()
	plaidCli := my_plaid.CreateClient()
	fbApp, _ := my_fb.CreateFirebaseApp()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowOrigins:     []string{
			"http://localhost:3000", 
			"http://localhost:3001", 
			"http://172.28.200.146:3000", 
			"https://trfer.me", 
			"https://goldfish-app-q2lpn.ondigitalocean.app",
		},
		// AllowOrigins:     []string{"http://172.28.38.241:3000", "http://172.28.38.241:3001"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST",},
		AllowHeaders:     []string{"Content-Type", "Authorization", "*"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.AuthMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Welcome to the usual api")
	})
	routes.CreateRoutes(router, plaidCli, psqlDB, s3Cli, fbApp)
	
	router.Run(":8080")
}