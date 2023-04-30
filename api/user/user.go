package user

import (
	"database/sql"
	"errors"

	"github.com/johnyeocx/usual/server2/api/auth"
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/models/user_models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	"github.com/johnyeocx/usual/server2/utils/enums/OtpType"
	"github.com/johnyeocx/usual/server2/utils/media"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)

func checkUsernameTaken(sqlDB *sql.DB, username string) (bool, *models.RequestError) {
	// step 1: check that email is not already taken
	u := user_db.UserDB{DB: sqlDB}
	_, err := u.GetUserByUsername(username)

	if err == sql.ErrNoRows {
		return false, nil;
	}

	if err != nil {
		return false, user_errors.GetUserFailedErr(err)
	}

	return true, nil
}

func checkEmailTaken(sqlDB *sql.DB, email string) (bool, *models.RequestError) {

	// step 1: check that email is not already taken
	u := user_db.UserDB{DB: sqlDB}
	_, err := u.GetUserByEmail(email)

	if err == sql.ErrNoRows {
		return false, nil;
	}

	if err != nil {
		return false, user_errors.GetUserFailedErr(err)
	}

	return true, nil
}

func getUserData(
	sqlDB *sql.DB, 
	uId int,
) (map[string]interface{}, *models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByID(uId)
	if err != nil {
		return nil, user_errors.GetUserFailedErr(err)
	}

	return map[string]interface{}{
		"user": user,
	}, nil
}

func getUser(
	sqlDB *sql.DB, 
	username string,
) (*user_models.User, *models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByUsername(username)
	if err != nil {
		return nil, user_errors.GetUserFailedErr(err)
	}

	return user, nil
}

func emailRegister(
	sqlDB *sql.DB, 
	username string, 
	email string,
) (*models.RequestError) {

	if (username == "" || email == "") {
		return user_errors.InvalidEmailErr(errors.New("invalid email"))
	}

	// step 1: check that username is not already taken
	taken, reqErr := checkUsernameTaken(sqlDB, email)
	if reqErr != nil {
		return reqErr
	}

	if taken {
		return user_errors.UsernameTakenErr(errors.New("username already taken"))
	}

	// step 2: check that email is not already taken
	taken, reqErr = checkEmailTaken(sqlDB, email)
	if reqErr != nil {
		return reqErr
	}

	if taken {
		return user_errors.EmailTakenErr(errors.New("email already taken"))
	}

	// step 3: create new entry
	u := user_db.UserDB{DB: sqlDB}
	_, err := u.CreateUserFromEmail(email, username)
	if err != nil {
		return user_errors.CreateUserFailedErr(err)
	}

	// step 4: send otp
	otp, reqErr := auth.GenerateEmailOtp(sqlDB, email, OtpType.EmailRegister)
	if reqErr != nil {
		return reqErr
	}
	
	_ = media.SendEmailVerification(email, *otp)
	return nil
}

func setName(
	sqlDB *sql.DB, 
	uId int,
	firstName string, 
	lastName string,
) (*models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	
	err := u.SetName(uId, firstName, lastName)
	
	if err != nil {
		return user_errors.SetNameFailedErr(err)
	}

	return nil
}

func initialiseBanking(
	sqlDB *sql.DB,
	plaidCli *plaid.APIClient,
	uId int,
	publicToken string,
) (*models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByID(uId)
	if err != nil {
		return user_errors.GetUserFailedErr(err)
	}

	// 1. Get access token
	accessToken, err := my_plaid.GetAuthAccessToken(plaidCli, publicToken)
	if err != nil {
		return banking_errors.GetAccessTokenFailedErr(err)
	}

	// 2. Get account bank details
	bacs, err := my_plaid.GetBACNumbers(plaidCli, accessToken)
	if err != nil {
		return banking_errors.GetBACSFailedErr(err)
	}
	
	recipientID, err := my_plaid.CreatePaymentRecipient(plaidCli, user.FullName(), bacs.Account, bacs.SortCode)
	if err != nil {
		return banking_errors.CreatePaymentRecipientFailedErr(err)
	}

	// 1. Set public token
	err = u.SetPublicTokenAndRecipientID(uId, publicToken, recipientID)
	if err != nil {
		return user_errors.SetPublicTokenFailedErr(err)
	}

	return nil
}