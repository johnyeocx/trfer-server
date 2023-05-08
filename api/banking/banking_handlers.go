package banking

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)

func Routes(bankingRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	bankingRouter.GET("/get_auth_link_token", getAuthLinkTokenHandler(sqlDB, plaidCli))
	bankingRouter.POST("/webhook", bankingWebhookHandler(sqlDB, plaidCli))
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


func bankingWebhookHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
	return func (c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		event := map[string]interface{}{}
		if err := json.Unmarshal(payload, &event); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}	

		webhookType := event["webhook_type"].(string)
		fmt.Println("Webhook type:", webhookType)
		fmt.Println("Webbhook event:", event)
		// switch webhookType {
		// 	case "PAYMENT_INITIATION": 
		// 		piEvent := decodePaymentInitiationWebhook(event)
		// 		fmt.Println("PIEvent:", piEvent)
		// 		reqErr := UpdatePaymentFromPIEvent(sqlDB, plaidCli, piEvent)
		// 		if reqErr != nil {
		// 			reqErr.LogAndReturn(c)
		// 			return
		// 		}
		// }

		c.JSON(200, nil)
	}
}