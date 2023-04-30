package auth_errors

import (
	"net/http"

	"github.com/johnyeocx/usual/server2/db/models"
)

type AuthError string
const (
	GenerateTokensFailed AuthError = "generate_tokens_failed"
	InvalidTokenType AuthError = "invalid_token_type"
	ParseRefreshTokenFailed AuthError = "parse_refresh_token_failed"
	ValidateUserFailed AuthError = "validate_user_failed"
)

func GenerateTokensFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(GenerateTokensFailed),
	}
}

func InvalidTokenTypeErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(InvalidTokenType),
	}
}

func ParseRefreshTokenFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusBadGateway,
		Code: string(ParseRefreshTokenFailed),
	}
}

func ValidateUserFailedErr(err error) *models.RequestError {
	return &models.RequestError{
		Err: err,
		StatusCode: http.StatusUnauthorized,
		Code: string(ValidateUserFailed),
	}
}