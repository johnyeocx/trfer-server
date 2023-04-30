package banking

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)

func Routes(bankingRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	bankingRouter.GET("/get_auth_link_token", getAuthLinkTokenHandler(sqlDB, plaidCli))
}

func getAuthLinkTokenHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}

		linkToken, err := my_plaid.CreateLinkToken(plaidCli, uId)
		if err != nil {
			reqErr := banking_errors.CreateAuthLinkTokenFailedErr(err)
			reqErr.LogAndReturn(c)
			return
		}
		
		c.JSON(http.StatusOK, map[string]string{
			"link_token": linkToken,
		})
	}
}
