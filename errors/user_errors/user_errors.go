package user_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type UserError string
const (
	GetUserFailed 		UserError = "get_user_failed"
	GetLastPaymentIDFailed 		UserError = "get_last_payment_id_failed"
	GetPaymentsFailed 		UserError = "get_payments_failed"
	NoEmailFound			UserError = "no_email_found"

	CreateUserFailed 	UserError = "create_user_failed"
	InvalidEmail 		UserError = "invalid_email"

	UsernameTaken 		UserError = "username_taken"
	EmailTaken 			UserError = "email_taken"
	SendEmailFailed 		UserError = "send_email_failed"

	SetEmailVerifiedFailed 		UserError = "set_email_verified_failed"
	SetNameFailed 		UserError = "set_name_failed"
	SetAccessTokenFailed 		UserError = "set_access_token_failed"
	SetRecipientIDFailed 		UserError = "set_recipient_id_failed"
	SetPageThemeFailed 		UserError = "set_page_theme_failed"
	SetUsernameFailed 		UserError = "set_username_failed"
	SetPersApprovedFailed 		UserError = "set_pers_approved_failed"
)

func GetUserFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetUserFailed),
	}
}

func GetLastPaymentIDFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetLastPaymentIDFailed),
	}
}

func GetPaymentsFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetPaymentsFailed),
	}
}


func NoEmailFoundErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusNotFound,
		Code: string(NoEmailFound),
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

func SetAccessTokenTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetAccessTokenFailed),
	}
}

func SetRecipientIDFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetRecipientIDFailed),
	}
}

func SetPageThemeFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetPageThemeFailed),
	}
}

func SetUsernameFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetUsernameFailed),
	}
}

func SetPersApprovedFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(SetPersApprovedFailed),
	}
}

