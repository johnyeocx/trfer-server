package pers_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type PersError string
const (
	CreateAccountFailed PersError = "create_pers_account_failed"
	DecodeInquiryFailed PersError = "decode_inquiry_failed"
	UpdateInquiryFailed PersError = "update_inquiry_failed"
)

func CreateAccountFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateAccountFailed),
	}
}

func DecodeInquiryFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(DecodeInquiryFailed),
	}
}


func UpdateInquiryFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(UpdateInquiryFailed),
	}
}