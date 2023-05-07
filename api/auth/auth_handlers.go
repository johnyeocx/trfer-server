package auth

import (
	"database/sql"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/utils/enums/TokenType"
	"github.com/johnyeocx/usual/server2/utils/helpers"
	"github.com/johnyeocx/usual/server2/utils/middleware"
)


func Routes(authRouter *gin.RouterGroup, sqlDB *sql.DB, fbApp *firebase.App) {
	authRouter.POST("/validate_token", validateTokenHandler(sqlDB))
	authRouter.POST("/refresh_user_token", refreshUserTokenHandler(sqlDB))
	
	authRouter.POST("/verify_email_register_otp", verifyEmailRegisterOtpHandler(sqlDB))

	// login
	authRouter.POST("/external_login", externalLoginHandler(sqlDB, fbApp))
	authRouter.POST("/login", loginHandler(sqlDB))
	authRouter.POST("/verify_email_login_otp", verifyEmailLoginOtpHandler(sqlDB))
}

func externalLoginHandler(sqlDB *sql.DB, fbApp *firebase.App) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		// 1. Get user email and search if exists in db
		reqBody := struct {
			Token  		 string `json:"token"`
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

		res, reqErr := externalLogin(sqlDB, email.(string))

		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return;
		}

		c.JSON(http.StatusOK, res)
	}
}


func loginHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		// 1. Get user email and search if exists in db
		reqBody := struct {
			Email  		 string `json:"email"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		reqErr := login(sqlDB, reqBody.Email)
		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, reqErr.Code)
			return
		}
		
		c.JSON(http.StatusOK, nil)
	}
}

func verifyEmailLoginOtpHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get user email and search if exists in db
		reqBody := struct {
			Email  	 	string `json:"email"`
			Otp  		 string `json:"otp"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		res, reqErr := verifyEmailLoginOtp(sqlDB, reqBody.Email, reqBody.Otp)
		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, reqErr.Code)
			return
		}
		// c.SetCookie("access_token", res["access_token"], 60 * 60 * 24, "/", "http://localhost", false, false);
		// c.SetCookie("refresh_token", res["refresh_token"], 60 * 60 * 24, "/", "http://localhost", false, false);

		c.JSON(http.StatusOK, res)
	}
}

func validateTokenHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		_, err := middleware.AuthenticateUser(c, sqlDB)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		c.JSON(200, nil)
	}
}

func refreshUserTokenHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {

		reqBody := struct {
			RefreshToken	string `json:"refresh_token"`
		}{}
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		res, reqErr := refreshToken(sqlDB, reqBody.RefreshToken, TokenType.User)

		if reqErr != nil {
			reqErr.LogAndReturn(c)
			return
		}

		// c.SetCookie("access_token", res["access_token"], 60 * 60 * 24, "/", "localhost", false, true);
		// c.SetCookie("refresh_token", res["refresh_token"], 60 * 60 * 24, "/", "localhost", false, true);
		c.JSON(http.StatusOK, res)
	}
}

func verifyEmailRegisterOtpHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get user email and search if exists in db
		reqBody := struct {
			Email  	 	string `json:"email"`
			Otp  		 string `json:"otp"`
		}{}
		
		if ok := helpers.DecodeReqBody(c, &reqBody); !ok {
			return
		}

		res, reqErr := verifyEmailRegisterOtp(sqlDB, reqBody.Email, reqBody.Otp)
		if reqErr != nil {
			reqErr.Log()
			c.JSON(reqErr.StatusCode, reqErr.Code)
			return
		}

		c.SetCookie("access_token", res["access_token"], 60 * 60 * 24, "/", "localhost", false, true);
		c.SetCookie("refresh_token", res["refresh_token"], 60 * 60 * 24, "/", "localhost", false, true);

		c.JSON(http.StatusOK, res)
	}
}