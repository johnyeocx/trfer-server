package transfer

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/plaid/plaid-go/v11/plaid"
)


func Routes(transferRouter *gin.RouterGroup, sqlDB *sql.DB, plaidCli *plaid.APIClient) {
	transferRouter.POST("/open_amount", transferOpenAmtHandler(sqlDB, plaidCli))
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