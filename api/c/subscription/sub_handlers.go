package subscription

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/db/models/sub_models"
	"github.com/johnyeocx/usual/server2/utils/enums"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/johnyeocx/usual/server2/utils/middleware"
)

func Routes(subRouter *gin.RouterGroup, sqlDB *sql.DB) {
	subRouter.POST("create", insertSubscriptionHandler(sqlDB))
}

// 1. Add custom
func insertSubscriptionHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		cusId, err := middleware.AuthenticateCId(c, sqlDB)
		if err != nil {
			return
		}
		
		reqBody := struct {
			BrandID		 	int 			`json:"brand_id"`
			Plan 			sub_models.Plan `json:"plan"`
			Status			enums.SubStatus	`json:"status"`
			}{}
			
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}
			

		subId, reqErr := CreateSub(sqlDB, *cusId, reqBody.BrandID, reqBody.Status, reqBody.Plan)
		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, reqErr.Code)
			return
		}
		
		c.JSON(http.StatusOK, subId)
	}
}

// 2. Detected

// 3. Custom