package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/otp_db"
	"github.com/johnyeocx/usual/server2/errors/otp_errors"
	"github.com/johnyeocx/usual/server2/utils/enums"
	"github.com/johnyeocx/usual/server2/utils/secure"
)

var (
    registerExpiry = time.Minute * 5
)


func GenerateEmailOtp(
    sqlDB *sql.DB, 
    email string,
    otpType enums.OtpType,
) (*string, *models.RequestError) {

	o := otp_db.OtpDB{DB: sqlDB}
	err := o.DeleteEmailOtp(email, otpType)

	if err != nil {
		return nil, otp_errors.DeleteOtpFailedErr(err)
	}

    // 2. generate otp
    otp := secure.GenerateOTP(6)
    hashedOtp, err := secure.GenerateHashFromStr(otp)
    if err != nil {
        return nil, otp_errors.HashOtpFailedErr(err)
    }

    // 3. insert verification into sql table
    expiry := time.Now().Add(registerExpiry).UTC()
    err = o.InsertEmailOtp(email, otpType, hashedOtp, expiry)
	if err != nil {
		return nil, otp_errors.InsertOtpFailedErr(err)
	}

    return &otp, nil
}

func VerifyEmailOtp(
    sqlDB *sql.DB, 
    email string,
    otp string,
    otpType enums.OtpType,
) (*models.RequestError) {
	o := otp_db.OtpDB{DB: sqlDB}

    // 1. Find matching verification in sql
    emailOtp, err := o.GetEmailOtp(email, otpType)

    if err != nil {
        return otp_errors.VerificationFailedErr(err)
    }

    // 2. Check otp match
    if !secure.StringMatchesHash(otp, emailOtp.HashedOTP) {
        return otp_errors.InvalidOtpErr(errors.New("otp given is invalid"))
    }

    // 3. Delete verification code from table
    if err := o.DeleteEmailOtp(email, otpType); err != nil {
        return otp_errors.DeleteOtpFailedErr(err)
    }

    return nil  
}
