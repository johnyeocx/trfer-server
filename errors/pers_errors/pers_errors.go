package pers_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type PersError string
const (
	CreateAccountFailed PersError = "create_pers_account_failed"
)

func CreateAccountFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateAccountFailed),
	}
}