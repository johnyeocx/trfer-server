package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	gen_errors "github.com/johnyeocx/usual/server2/errors/general_errors"
	"github.com/johnyeocx/usual/server2/errors/util_errors"
	"github.com/johnyeocx/usual/server2/utils/cloud"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	"github.com/plaid/plaid-go/v11/plaid"
)

func Routes(userRouter *gin.RouterGroup, sqlDB *sql.DB, s3Cli *s3.Client, plaidCli *plaid.APIClient) {
	userRouter.POST("/email_register", emailRegisterHandler(sqlDB))
	userRouter.POST("/name_and_photo", setNameAndPhotoHandler(sqlDB, s3Cli))
	userRouter.POST("/initialise_banking", initialiseBankingHandler(sqlDB, plaidCli))
	userRouter.GET("/data", getUserDataHandler(sqlDB))
	userRouter.GET("/:username", getUserHandler(sqlDB))
}

func emailRegisterHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get user email and search if exists in db
		reqBody := struct {
			Username  	 string `json:"username"`
			Email  		 string `json:"email"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		reqErr := emailRegister(sqlDB, reqBody.Username, reqBody.Email)
		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, reqErr.Code)
			return
		}
		
		c.JSON(http.StatusOK, nil)
	}
}

func setNameAndPhotoHandler(sqlDB *sql.DB, s3Cli *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			FirstName  	 string `json:"first_name"`
			LastName  		 string `json:"last_name"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		// 2. Set Name
		reqErr := setName(sqlDB, uId, reqBody.FirstName, reqBody.LastName)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		// 3. Generate Upload Link
		key := fmt.Sprintf("user/profile_image/%d", uId)
		err = cloud.DeleteObject(s3Cli, key)
		uploadUrl, err := cloud.GetImageUploadUrl(s3Cli, key)
		if err != nil {
			reqErr := util_errors.GetPresignedUrlFailedErr(err)
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, map[string]string {
			"upload_url": uploadUrl,
		})
	}
}

func initialiseBankingHandler(sqlDB *sql.DB, plaidCli *plaid.APIClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			PublicToken  	 string `json:"public_token"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		// 2. Set Name
		reqErr := initialiseBanking(sqlDB, plaidCli, uId, reqBody.PublicToken)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func getUserDataHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}

		// 2. Set Name
		res, reqErr := getUserData(sqlDB, uId)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func getUserHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		username, ok := c.Params.Get("username")

		if !ok  {
			reqErr := gen_errors.InvalidRequestParamErr(errors.New("username not valid"))
			reqErr.LogAndReturn(c)
			return;
		}

		// 2. Set Name
		user, reqErr := getUser(sqlDB, username)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}