package c_auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/johnyeocx/usual/server2/db/cus_db"
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/errors/cus_errors"
	"github.com/johnyeocx/usual/server2/utils/enums"
	"github.com/johnyeocx/usual/server2/utils/secure"
)


func ExternalSignIn(
	sqlDB *sql.DB,
	email string,
	signInProvider enums.SignInProvider,
) (map[string]interface{}, *models.RequestError) {
	c := cus_db.CustomerDB{DB: sqlDB}
	
	// 1. Check if email already exists
	cus, err := c.GetCustomerByEmail(email)
	if err != nil && err != sql.ErrNoRows{
		return nil, cus_errors.GetCusFailedErr(err)
	}

	// 2. If email doesn't exist, create new entry
	var cusId int
	if cus == nil || !cus.EmailVerified {
		cId, err := c.CreateCustomerFromExtSignin(email, signInProvider)
		if err != nil {
			return nil, cus_errors.CreateCusFailedErr(err)
		}
		cusId = *cId
	} else {
		cusId = cus.ID
	}

	accessToken, refreshToken, err := secure.GenerateTokensFromId(cusId, enums.CusUserType)
	if err != nil {
		return nil, &models.RequestError{
			Err: err,
			StatusCode: http.StatusBadGateway,
		}
	}

	return map[string]interface{}{
		"access_token": *accessToken,
		"refresh_token": *refreshToken,
	}, nil
}

func refreshToken(sqlDB *sql.DB, refreshToken string) (*int, error) {
	cusId, cusType, err := secure.ParseRefreshToken(refreshToken)
	if cusType != string(enums.CusUserType) {
		return nil, errors.New("unauthorized user")
	}

	if err != nil {
		return nil, err
	}
	
	customerIdInt, err := strconv.Atoi(cusId)
	if err != nil {
		return nil, err
	}
	
	if ok := cus_db.ValidateCustomerId(sqlDB, customerIdInt); !ok {
		return nil, fmt.Errorf("invalid business id")
	}

	return &customerIdInt, nil
}