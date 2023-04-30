package otp_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type OtpError string
const (
	DeleteOtpFailed OtpError = "delete_otp_failed"
	InsertOtpFailed OtpError = "insert_otp_failed"
	HashOtpFailed 	OtpError = "hash_otp_failed"
	VerificationFailed 	OtpError = "verification_failed"
	InvalidOtp 	OtpError = "InvalidOtp"
)

func InsertOtpFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(InsertOtpFailed),
	}
}


func DeleteOtpFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(DeleteOtpFailed),
	}
}

func HashOtpFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(HashOtpFailed),
	}
}

func VerificationFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(VerificationFailed),
	}
}

func InvalidOtpErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusForbidden,
		Code: string(InvalidOtp),
	}
}