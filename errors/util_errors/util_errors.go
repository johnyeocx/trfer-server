package util_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type UtilError string
const (
	GetPresignedUrlFailed UtilError = "get_presigned_url_failed"
)

func GetPresignedUrlFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GetPresignedUrlFailed),
	}
}