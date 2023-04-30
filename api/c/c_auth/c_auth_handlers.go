package c_auth

import (
	"database/sql"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/utils/enums"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	"github.com/johnyeocx/usual/server2/utils/secure"
)

// AUTH ROUTES
func Routes(authRouter *gin.RouterGroup, sqlDB *sql.DB, fbApp *firebase.App) {

	authRouter.POST("/external_sign_in", externalSignInHandler(sqlDB, fbApp))
	authRouter.POST("/validate", validateTokenHandler(sqlDB))
	authRouter.POST("/refresh_token", refreshTokenHandler(sqlDB))
	// authRouter.POST("/login", loginHandler(sqlDB))
}

func externalSignInHandler(sqlDB *sql.DB, fbApp *firebase.App) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			Token  		 string `json:"token"`
		}{}

		if err := c.BindJSON(&reqBody); err != nil {
			log.Println("failed to decode req body: ", err)
			c.JSON(400, err)
			return
		}

		client, err := fbApp.Auth(c)
		if err != nil {
			log.Println("failed to get auth client: ", err)
			c.JSON(400, err)
			return
		}

		// Verify the ID token first.
		token, err := client.VerifyIDToken(c, reqBody.Token)
		if err != nil {
			log.Println("failed to get verify token: ", err)
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		email := token.Claims["email"]
		siginProvider := token.Firebase.SignInProvider
		res, reqErr := ExternalSignIn(sqlDB, email.(string), enums.SignInProvider(siginProvider))

		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, res)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func validateTokenHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {

		_, err := middleware.AuthenticateCId(c, sqlDB)

		if err != nil {
			log.Println("Failed to validate customer:", err)
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		c.JSON(200, nil)
	}
}

func refreshTokenHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {

		refreshTokenStr, err := c.Cookie("refresh_token")
		if err != nil || len(refreshTokenStr) == 0 {
			reqBody := struct {
				RefreshToken	string `json:"refresh_token"`
			}{}
	
			if err := c.BindJSON(&reqBody); err != nil {
				log.Printf("Failed to decode req body for refresh token: %v\n", err)
				c.JSON(400, err)
				return
			}
			refreshTokenStr = reqBody.RefreshToken
		}

		cId, err := refreshToken(sqlDB, refreshTokenStr)
		if err != nil {
			log.Printf("Failed to authenticate refresh token: %v\n", err)
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		accessToken, refreshToken, err := secure.GenerateTokensFromId(*cId, "customer")
		if err != nil {
			c.JSON(http.StatusBadGateway, err)
			return
		}

		c.JSON(200, map[string]string{
			"access_token": *accessToken,
			"refresh_token": *refreshToken,
		})
	}
}

// func loginHandler(sqlDB *sql.DB) gin.HandlerFunc {
// 	return func (c *gin.Context) {
		
// 		reqBody := struct {
// 			Email  		 string `json:"email"`
// 			Password     string `json:"password"`
// 		}{}

// 		if err := c.BindJSON(&reqBody); err != nil {
// 			log.Printf("Failed to decode req body for login: %v\n", err)
// 			c.JSON(400, err)
// 			return
// 		}

// 		res, reqErr := login(sqlDB, reqBody.Email, reqBody.Password)

// 		if reqErr != nil {
// 			log.Println("Failed to login:", reqErr.Err)
// 			c.JSON(reqErr.StatusCode, res)
// 			return 
// 		}
		
// 		accessToken := res["access_token"].(string)
// 		refreshToken := res["refresh_token"].(string)

// 		c.SetCookie("access_token", accessToken, 60 * 60 * 24, "/", "localhost", false, true);
// 		c.SetCookie("refresh_token", refreshToken, 60 * 60 * 24, "/", "localhost", false, true);

// 		c.JSON(http.StatusOK, res)
// 	}
// }

