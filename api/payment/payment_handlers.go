package payment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/utils/enums/PaymentStatus"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/plaid/plaid-go/v11/plaid"
)


func Routes(transferRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	transferRouter.POST("/open_amount", transferOpenAmtHandler(sqlDB, plaidCli))
	transferRouter.POST("/webhook", transferWebhookHandler(sqlDB, plaidCli))
}

func transferOpenAmtHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
	return func (c *gin.Context) {
		// 1. Get user email and search if exists in db
		reqBody := struct {
			ToUsername  string 		`json:"to_username"`
			Amount		float64		`json:"amount"`
			Note		string		`json:"note"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		if (reqBody.Amount > 5.00) {
			fmt.Println("Amt too high!")
			c.JSON(http.StatusBadRequest, errors.New("amount too high"))
			return;
		}

		linkToken, reqErr := TransferOpenAmt(sqlDB, plaidCli, reqBody.ToUsername, reqBody.Amount, reqBody.Note)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}
		
		c.JSON(http.StatusOK, map[string]string {
			"link_token": linkToken,
		})
	}
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

		fmt.Println("Webhook event:", event)
		webhookType := event["webhook_type"].(string)
		fmt.Println("Type:", webhookType)
		switch webhookType {
			case "PAYMENT_INITIATION": 
				piEvent := decodePaymentInitiationWebhook(event)
				fmt.Println("PIEvent:", piEvent)
				reqErr := UpdatePaymentFromPIEvent(sqlDB, plaidCli, piEvent)
				if reqErr != nil {
					reqErr.LogAndReturn(c)
					return
				}
				fmt.Println("Successfully inserted")
		}

		c.JSON(200, nil)
	}
}

func decodePaymentInitiationWebhook (event map[string]interface{}) (models.PaymentInitiationEvent) {
	
	paymentStatus := event["new_payment_status"].(string)
	paymentId := event["payment_id"].(string)

	return models.PaymentInitiationEvent{
		NewPaymentStatus: PaymentStatus.EventStrToPaymentStatus(paymentStatus),
		PlaidPaymentID: paymentId,
	}
}
