package banking_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type BankingError string
const (
	CreateAuthLinkTokenFailed BankingError = "create_auth_link_token_failed"
	CreatePaymentLinkTokenFailed BankingError = "create_payment_link_token_failed"
	GetAccessTokenFailed BankingError = "get_access_token_failed"
	GetBACSFailed BankingError = "get_bacs_failed"
	CreatePaymentRecipientFailed BankingError = "create_payment_recipient_failed"
	CreatePaymentFailed BankingError = "create_payment_failed"

	InsertPaymentFailed BankingError = "insert_payment_failed"
	UpdatePaymentFailed BankingError = "update_payment_failed"
	GetPaymentFailed BankingError = "get_payment_failed"
	GetTransactionsFailed BankingError = "get_transactions_failed"

	AlreadySetAccessToken BankingError = "already_set_access_token"
	NoAccessToken BankingError = "no_access_token"
	AlreadyCreatedRecipientID BankingError = "already_created_recipient_id"
)

func CreateAuthLinkTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateAuthLinkTokenFailed),
	}
}

func CreatePaymentLinkTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreatePaymentLinkTokenFailed),
	}
}

func GetAccessTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetAccessTokenFailed),
	}
}

func GetPaymentFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetPaymentFailed),
	}
}

func GetTransactionsFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetTransactionsFailed),
	}
}


func GetBACSFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetBACSFailed),
	}
}

func CreatePaymentRecipientFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreatePaymentRecipientFailed),
	}
}

func CreatePaymentFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreatePaymentFailed),
	}
}

func InsertPaymentFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(InsertPaymentFailed),
	}
}

func UpdatePaymentFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(UpdatePaymentFailed),
	}
}

func AlreadySetAccessTokenErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(AlreadySetAccessToken),
	}
}

func NoAccessTokenErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(AlreadySetAccessToken),
	}
}

func AlreadyCreatedRecipientIDErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(AlreadyCreatedRecipientID),
	}
}

