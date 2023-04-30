package sub_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type SubError string
const (
	CreateSubFailed SubError = "create_sub_failed"
	InvalidBrandID	SubError = "invalid_brand_id"
	InvalidCusID	SubError = "invalid_cus_id"
	DuplicateSub	SubError = "duplicate_sub"
)


func InvalidCusIDErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadRequest,
		Code: string(InvalidCusID),
	}
}

func InvalidBrandIDErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadRequest,
		Code: string(InvalidBrandID),
	}
}

func DuplicateSubErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusConflict,
		Code: string(DuplicateSub),
	}
}

func CreateSubFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(CreateSubFailed),
	}
}