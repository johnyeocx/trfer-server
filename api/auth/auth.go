package auth

import (
	"database/sql"
	"strconv"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/auth_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	"github.com/johnyeocx/usual/server2/utils/enums"
	"github.com/johnyeocx/usual/server2/utils/enums/OtpType"
	"github.com/johnyeocx/usual/server2/utils/enums/TokenType"
	"github.com/johnyeocx/usual/server2/utils/media"
	"github.com/johnyeocx/usual/server2/utils/secure"
)

func login (sqlDB *sql.DB, email string) (*models.RequestError) {
	u := user_db.UserDB{DB: sqlDB}

	_, err := u.GetUserByEmail(email)

	if err != nil {
		return user_errors.GetUserFailedErr(err)
	}

		// step 4: send otp
	otp, reqErr := GenerateEmailOtp(sqlDB, email, OtpType.EmailLogin)
	if reqErr != nil {
		return reqErr
	}
	
	_ = media.SendEmailVerification(email, *otp)
	return nil
}

func verifyEmailLoginOtp (sqlDB *sql.DB, email string, otp string) (map[string]string, *models.RequestError){
	reqErr := VerifyEmailOtp(sqlDB, email, otp, OtpType.EmailLogin)
	if reqErr != nil {
		return nil, reqErr
	}

	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return nil, user_errors.GetUserFailedErr(err)
	}
	
	accessToken, refreshToken, err := secure.GenerateTokenPair(user.ID, TokenType.User)
	if err != nil {
		return nil, auth_errors.GenerateTokensFailedErr(err)
	}
	
	return map[string]string {
		"access_token": accessToken,
		"refresh_token": refreshToken,
	}, nil
}


func verifyEmailRegisterOtp (sqlDB *sql.DB, email string, otp string) (map[string]string, *models.RequestError){

	// 1. Verify Email
	reqErr := VerifyEmailOtp(sqlDB, email, otp, OtpType.EmailRegister)
	if reqErr != nil {
		return nil, reqErr
	}

	// 2. Set User Verified
	u := user_db.UserDB{DB: sqlDB}
	user, err := u.SetEmailVerified(email)
	if err != nil {
		return nil, user_errors.SetEmailVerifiedFailedErr(err)
	}

	accessToken, refreshToken, err := secure.GenerateTokenPair(user.ID, TokenType.User)
	if err != nil {
		return nil, auth_errors.GenerateTokensFailedErr(err)
	}
	
	return map[string]string {
		"access_token": accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func refreshToken(sqlDB *sql.DB, refreshToken string, tokenType enums.TokenType) (map[string]string, *models.RequestError) {
	uId, tokenType, err := secure.ParseRefreshToken(refreshToken)

	if err != nil {
		return nil, auth_errors.ParseRefreshTokenFailedErr(err)
	}

	if tokenType != TokenType.User {
		return nil, auth_errors.InvalidTokenTypeErr(err)
	}
	
	userIdInt, _ := strconv.Atoi(uId)
	
	if ok := user_db.ValidateUser(sqlDB, userIdInt); !ok {
		return nil, auth_errors.ValidateUserFailedErr(err)
	}

	accessToken, refreshToken, err := secure.GenerateTokenPair(userIdInt, tokenType)
	if err != nil {
		return nil, auth_errors.GenerateTokensFailedErr(err)
	}

	return map[string]string{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	}, nil
}