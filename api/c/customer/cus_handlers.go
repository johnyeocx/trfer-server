package customer

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Routes(customerRouter *gin.RouterGroup, sqlDB *sql.DB) {
	// customerRouter.POST("create", createCustomerHandler(sqlDB))
	// customerRouter.POST("verify_email", verifyCustomerEmailHandler(sqlDB, s3Sess))
	// customerRouter.POST("resend_email_otp", resendEmailOTPHandler(sqlDB))

	// customerRouter.PATCH("name", updateCusNameHandler(sqlDB))
	// customerRouter.PATCH("email", sendCusUpdateEmailVerificationHandler(sqlDB))
	// customerRouter.PATCH("verify_email", verifyCusUpdateEmailHandler(sqlDB))
	// customerRouter.PATCH("address", updateCusAddressHandler(sqlDB))
	// customerRouter.PATCH("password", updateCusPasswordHandler(sqlDB))
	// customerRouter.PATCH("default_payment", updateCusDefaultPaymentHandler(sqlDB))
}



// func createCustomerHandler(sqlDB *sql.DB) gin.HandlerFunc {
// 	return func (c *gin.Context) {
// 		reqBody := struct {
// 			FirstName			string 	`json:"first_name"`
// 			LastName			string 	`json:"last_name"`
// 			Email 			string 	`json:"email"`
// 			Password	 	string 	`json:"password"`
// 		}{}

// 		if err := c.BindJSON(&reqBody); err != nil {
// 			log.Printf("Failed to decode req body for verify otp: %v\n", err)
// 			c.JSON(400, err)
// 			return
// 		}

// 		reqErr := CreateCustomer(sqlDB, reqBody.FirstName, reqBody.LastName, reqBody.Email, reqBody.Password)
// 		if reqErr != nil {
// 			log.Println("Failed to create customer:", reqErr.Err)
// 			c.JSON(reqErr.StatusCode, reqErr.Err)
// 			return
// 		}

// 		c.JSON(200, nil)
// 	}
// }

// func verifyCustomerEmailHandler(sqlDB *sql.DB, s3Sess *session.Session) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		// 1. Get user email and search if exists in db
// 		reqBody := struct {
// 			Email  		 string `json:"email"`
// 			OTP          string `json:"otp"`
// 		}{}

		
// 		if err := c.BindJSON(&reqBody); err != nil {
// 			log.Printf("Failed to decode req body for verify otp: %v\n", err)
// 			c.JSON(400, err)
// 			return
// 		}
// 		fmt.Println(reqBody)
		
// 		// 2. Verify email
// 		res, reqErr := VerifyCustomerRegEmail(s3Sess, sqlDB, reqBody.Email, reqBody.OTP)
// 		if reqErr != nil {
// 			log.Println(reqErr.Err)
// 			c.JSON(reqErr.StatusCode, reqErr.Err)
// 			return
// 		}

// 		c.JSON(200, res)
// 	}
// }

// func resendEmailOTPHandler(sqlDB *sql.DB) gin.HandlerFunc {
// 	return func (c* gin.Context) {
// 		reqBody := struct {
// 			Email 	string	`json:"email"`
// 			OtpType  string  `json:"otp_type"`
// 		}{}

// 		if err := c.BindJSON(&reqBody); err != nil {
// 			log.Printf("Failed to decode req body: %v\n", err)
// 			c.JSON(400, err)
// 			return
// 		}
		
// 		var reqErr *models.RequestError
// 		if reqBody.OtpType == constants.OtpTypes.RegisterCusEmail  {

// 			reqErr = sendRegEmailOTP(sqlDB, reqBody.Email)
// 		} else if reqBody.OtpType == constants.OtpTypes.UpdateCusEmail {
// 			cusId, err := middleware.AuthenticateCId(c, sqlDB)
// 			if err != nil {
// 				c.JSON(http.StatusUnauthorized, err)
// 				return
// 			}
// 			reqErr = sendUpdateEmailOTP(sqlDB, *cusId, reqBody.Email)
// 		} else {
// 			c.JSON(http.StatusBadRequest, errors.New("invalid otp type"))
// 		}

// 		if reqErr != nil {
// 			log.Printf("Failed to resend cus update email verification: %v\n", reqErr.Err)
// 			c.JSON(http.StatusBadGateway, reqErr.Err)
// 			return
// 		}

// 		c.JSON(200, nil)
// 	}
// }

// func addCustomerCardHandler(sqlDB *sql.DB) gin.HandlerFunc {
// 	return func (c *gin.Context) {
// 		cusId, err := middleware.AuthenticateCId(c, sqlDB)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, err)
// 		}

// 		reqBody := struct {
// 			Number		string `json:"number"`
// 			ExpMonth 	int64 `json:"expiry_month"`
// 			ExpYear 	int64 `json:"expiry_year"`
// 			CVC 		string `json:"cvc"`
// 		}{}

// 		if err := c.BindJSON(&reqBody); err != nil {
// 			log.Printf("Failed to decode req body: %v\n", err)
// 			c.JSON(400, err)
// 			return
// 		}

// 		res, reqErr := AddCusCreditCard(sqlDB, *cusId, models.CreditCard{
// 			Number: reqBody.Number,
// 			ExpMonth: reqBody.ExpMonth,
// 			ExpYear: reqBody.ExpYear,
// 			CVC: reqBody.CVC,
// 		})

// 		if reqErr != nil {
// 			log.Println("Failed to add custoemr credit card:", reqErr.Err)
// 			c.JSON(reqErr.StatusCode, reqErr.Err)
// 			return
// 		}

// 		c.JSON(200, res)
// 	}
// }
