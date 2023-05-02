package transfer

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/plaid/plaid-go/v11/plaid"
)


func Routes(transferRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	transferRouter.POST("/open_amount", transferOpenAmtHandler(sqlDB, plaidCli))
	transferRouter.POST("/webhook", transferWebhookHandler(sqlDB))
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
			fmt.Println("Amt too hight")
			c.JSON(http.StatusBadRequest, errors.New("amount too high"))
		}

		linkToken, reqErr := TransferOpenAmt(sqlDB, plaidCli, reqBody.ToUsername, reqBody.Amount)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}
		
		c.JSON(http.StatusOK, map[string]string {
			"link_token": linkToken,
		})
	}
}

func transferWebhookHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		fmt.Println("Got request body")
		payload, err := io.ReadAll(c.Request.Body)
		fmt.Println("Read in req body:", payload)
		if err != nil {
			fmt.Println("Failed to io.readall:", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		
		event := map[string]interface{}{}
		if err := json.Unmarshal(payload, &event); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}	

		fmt.Println("EVENT:", event)

		c.JSON(200, nil)
	}
}
