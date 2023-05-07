package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	gen_errors "github.com/johnyeocx/usual/server2/errors/general_errors"
	"github.com/johnyeocx/usual/server2/errors/util_errors"
	"github.com/johnyeocx/usual/server2/utils/cloud"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/johnyeocx/usual/server2/utils/middleware"
	"github.com/plaid/plaid-go/v11/plaid"
)

func Routes(userRouter *gin.RouterGroup, sqlDB *sql.DB, s3Cli *s3.Client, plaidCli *plaid.APIClient, fbApp *firebase.App) {
	userRouter.POST("/external_register", externalRegisterHandler(sqlDB, fbApp))
	userRouter.POST("/email_register", emailRegisterHandler(sqlDB))
	userRouter.POST("/name_and_photo", setNameAndPhotoHandler(sqlDB, s3Cli))
	userRouter.POST("/initialise_banking", initialiseBankingHandler(sqlDB, plaidCli))

	userRouter.PATCH("/profile_image", setProfileImgHandler(sqlDB, s3Cli))
	userRouter.PATCH("/name", setNameHandler(sqlDB))
	userRouter.PATCH("/username", setUsernameHandler(sqlDB))
	userRouter.PATCH("/page_theme", setPageThemeHandler(sqlDB))

	userRouter.GET("/data", getUserDataHandler(sqlDB))
	userRouter.GET("/:username", getUserHandler(sqlDB))
}


func externalRegisterHandler(sqlDB *sql.DB, fbApp *firebase.App) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			Token  		 string `json:"token"`
			Username	string 	`json:"username"`
		}{}

		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		client, err := fbApp.Auth(c)
		if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
		}

		// Verify the ID token first.
		token, err := client.VerifyIDToken(c, reqBody.Token)
		if err != nil {
				log.Fatal(err)
		}

		email := token.Claims["email"]

		res, reqErr := ExternalRegister(sqlDB, email.(string), reqBody.Username)

		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return;
		}

		c.JSON(http.StatusOK, res)
	}
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

func setProfileImgHandler(sqlDB *sql.DB, s3Cli *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 3. Generate Upload Link
		key := fmt.Sprintf("user/profile_image/%d", uId)
		cloud.DeleteObject(s3Cli, key)
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

func setNameHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			FirstName  	 string `json:"first_name"`
			LastName  	 string `json:"last_name"`
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

		c.JSON(http.StatusOK, nil)
	}
}

func setUsernameHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			Username  	 string `json:"username"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		// 2. Set Name
		reqErr := setUsername(sqlDB, uId, reqBody.Username)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func setPageThemeHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		uId, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			return
		}
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			PageTheme  	 string `json:"page_theme"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		// 2. Set Name
		reqErr := setPageTheme(sqlDB, uId, reqBody.PageTheme)
		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		c.JSON(http.StatusOK, nil)
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