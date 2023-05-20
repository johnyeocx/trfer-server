package persona

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/v11/plaid"
)


func Routes(persRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	persRouter.POST("/webhook", transferWebhookHandler(sqlDB, plaidCli))
}

func transferWebhookHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
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

		attributes := event["attributes"].(map[string]interface{})
		webhookType := attributes["name"].(string)

		fmt.Println("Received Pers Webhook: ", webhookType)
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
