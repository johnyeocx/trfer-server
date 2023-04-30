package user_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type UserError string
const (
	GetUserFailed 		UserError = "get_user_failed"
	CreateUserFailed 	UserError = "create_user_failed"
	InvalidEmail 		UserError = "invalid_email"

	UsernameTaken 		UserError = "username_taken"
	EmailTaken 			UserError = "email_taken"
	SendEmailFailed 		UserError = "send_email_failed"

	SetEmailVerifiedFailed 		UserError = "set_email_verified_failed"
	SetNameFailed 		UserError = "set_name_failed"
	SetPublicTokenFailed 		UserError = "set_public_token_failed"
)

func GetUserFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetUserFailed),
	}
}

func CreateUserFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateUserFailed),
	}
}

func InvalidEmailErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadRequest,
		Code: string(InvalidEmail),
	}
}

func UsernameTakenErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusConflict,
		Code: string(UsernameTaken),
	}
}

func EmailTakenErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusConflict,
		Code: string(EmailTaken),
	}
}

func SendEmailFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusConflict,
		Code: string(SendEmailFailed),
	}
}

func SetEmailVerifiedFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetEmailVerifiedFailed),
	}
}

func SetNameFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetNameFailed),
	}
}

func SetPublicTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetPublicTokenFailed),
	}
}
