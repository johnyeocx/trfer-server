package routes

import (
	"database/sql"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/api/auth"
	"github.com/johnyeocx/usual/server2/api/banking"
	"github.com/johnyeocx/usual/server2/api/payment"
	persona "github.com/johnyeocx/usual/server2/api/pers"
	"github.com/johnyeocx/usual/server2/api/user"
	"github.com/plaid/plaid-go/v11/plaid"
)



func CreateRoutes(
	router *gin.Engine,
	plaidCli *plaid.APIClient,
	sqlDB *sql.DB,
	s3Cli *s3.Client,
	fbApp *firebase.App,
) {
	apiRoute := router.Group("/api")
	{
		// c_auth.Routes(apiRoute.Group("/c/auth"), sqlDB, fbApp)
		// subscription.Routes(apiRoute.Group("/c/sub"), sqlDB)
		user.Routes(apiRoute.Group("/user"), sqlDB, s3Cli, plaidCli, fbApp)
		auth.Routes(apiRoute.Group("/auth"), sqlDB, fbApp)
		banking.Routes(apiRoute.Group("/banking"), sqlDB, plaidCli)
		payment.Routes(apiRoute.Group("/payment"), sqlDB, plaidCli)
		persona.Routes(apiRoute.Group("/pers"), sqlDB, plaidCli)
	}
}
