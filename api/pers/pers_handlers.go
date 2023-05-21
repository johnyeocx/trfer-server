package persona

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	"github.com/plaid/plaid-go/v11/plaid"
)


func Routes(persRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	persRouter.GET("/session_token", getInquiryAccessTokenHandler(sqlDB))
	persRouter.POST("/inquiry/webhook", persInquiryWebhookHandler(sqlDB, plaidCli))
}

func getInquiryAccessTokenHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}

		// 2. Set Name
		res, reqErr := GetInquiryAccessToken(sqlDB, uId)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func persInquiryWebhookHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
	return func (c *gin.Context) {
	

		const MaxBodyBytes = int64(200000)
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
		
		data := event["data"].(map[string]interface{})
	
		attributes := data["attributes"].(map[string]interface{})

		webhookName := attributes["name"].(string)
		fmt.Println("Received Pers Webhook: ", webhookName)

		attrData := attributes["payload"].(map[string]interface{})["data"].(map[string]interface{})

		webhookType := attrData["type"]
		switch webhookType {
			case "inquiry": 
				inquiry, reqErr := DecodeInquiryWebhook(attrData)
				if reqErr != nil {
					reqErr.LogAndReturn(c)
					return
				}
				reqErr = UpdateInquiry(sqlDB, *inquiry)
				if reqErr != nil {
					reqErr.LogAndReturn(c)
					return
				}
			case "inquiry-session":
				inquiry, reqErr := DecodeInquirySessionWebhook(attrData)
				if reqErr != nil {
					reqErr.LogAndReturn(c)
					return
				}
				reqErr = UpdateInquirySession(sqlDB, *inquiry)
				if reqErr != nil {
					reqErr.LogAndReturn(c)
					return
				}
		}

		c.JSON(200, nil)
	}
}