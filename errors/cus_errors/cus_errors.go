package cus_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type CusError string
const (
	GetCusFailed 	CusError = "get_customer_failed"
	CreateCusFailed CusError = "create_customer_failed"
)

func GetCusFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetCusFailed),
	}
}

func CreateCusFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateCusFailed),
	}
}