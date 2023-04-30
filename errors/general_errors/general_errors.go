package gen_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type GenError string
const (
	InvalidRequestParam 	GenError = "invalid_request_param"
)

func InvalidRequestParamErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadRequest,
		Code: string(InvalidRequestParam),
	}
}