package user

import (
	"database/sql"
	"errors"

	"github.com/johnyeocx/usual/server2/api/auth"
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/models/user_models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/auth_errors"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	gen_errors "github.com/johnyeocx/usual/server2/errors/general_errors"
	"github.com/johnyeocx/usual/server2/errors/pers_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	"github.com/johnyeocx/usual/server2/persona"
	"github.com/johnyeocx/usual/server2/utils/enums/OtpType"
	"github.com/johnyeocx/usual/server2/utils/enums/PaymentStatus"
	"github.com/johnyeocx/usual/server2/utils/enums/TokenType"
	"github.com/johnyeocx/usual/server2/utils/media"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/johnyeocx/usual/server2/utils/secure"
	"github.com/plaid/plaid-go/v11/plaid"
)

func CheckUsernameTaken(sqlDB *sql.DB, username string) (bool, *models.RequestError) {
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

func CheckEmailTaken(sqlDB *sql.DB, email string) (bool, *models.RequestError) {

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

	if (user.AccessToken.Valid) {
		user.AccessToken.Valid = false
		user.AccessToken.String = ""
		user.AccessTokenCreated = true
	} else {
		user.AccessTokenCreated = false
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

func getUserPayments(
	sqlDB *sql.DB, 
	uId int,
) ([]models.Payment, *models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	payments, err := u.GetUserPaymentsByStatus(uId, PaymentStatus.Executed)

	if err != nil {
		return nil, user_errors.GetUserFailedErr(err)
	}

	return payments, nil
}

func ExternalRegister(
	sqlDB *sql.DB, 
	email string, 
	username string,
) (map[string]string, *models.RequestError) {
	if (username == "" || email == "") {
		return nil, user_errors.InvalidEmailErr(errors.New("invalid email"))
	}

	// step 1: check that username is not already taken
	taken, reqErr := CheckUsernameTaken(sqlDB, username)
	if reqErr != nil {
		return nil, reqErr
	}

	if taken {
		return nil, user_errors.UsernameTakenErr(errors.New("username already taken"))
	}

	// step 2: check that email is not already taken
	taken, reqErr = CheckEmailTaken(sqlDB, email)
	if reqErr != nil {
		return nil, reqErr
	}

	if taken {
		return nil, user_errors.EmailTakenErr(errors.New("email already taken"))
	}

	// 2. Create pers acct
	persId, err := persona.CreateAccount(email)
	if err != nil {
		return nil, pers_errors.CreateAccountFailedErr(err)
	}

	// step 3: create new entry
	u := user_db.UserDB{DB: sqlDB}
	uId, err := u.CreateUserFromEmail(email, username, true, &persId)
	if err != nil {
		return nil, user_errors.CreateUserFailedErr(err)
	}

	accessToken, refreshToken, err := secure.GenerateTokenPair(*uId, TokenType.User)
	if err != nil {
		return nil, auth_errors.GenerateTokensFailedErr(err)
	}

	return map[string]string {
		"access_token": accessToken,
		"refresh_token": refreshToken,
	}, nil
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
	taken, reqErr := CheckUsernameTaken(sqlDB, email)
	if reqErr != nil {
		return reqErr
	}

	if taken {
		return user_errors.UsernameTakenErr(errors.New("username already taken"))
	}

	// step 2: check that email is not already taken
	taken, reqErr = CheckEmailTaken(sqlDB, email)
	if reqErr != nil {
		return reqErr
	}

	if taken {
		return user_errors.EmailTakenErr(errors.New("email already taken"))
	}

	// step 3: create new entry
	u := user_db.UserDB{DB: sqlDB}
	_, err := u.CreateUserFromEmail(email, username, false, nil)
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
	if (firstName == "" || lastName == "") {
		return gen_errors.InvalidRequestParamErr(errors.New("first or last name is empty"))
	}

	u := user_db.UserDB{DB: sqlDB}
	
	err := u.SetName(uId, firstName, lastName)
	
	if err != nil {
		return user_errors.SetNameFailedErr(err)
	}

	return nil
}

func setUsername(
	sqlDB *sql.DB, 
	uId int,
	username string, 
) (*models.RequestError) {

	if (username == "") {
		return gen_errors.InvalidRequestParamErr(errors.New("username is empty"))
	}

	u := user_db.UserDB{DB: sqlDB}
	taken, reqErr := CheckUsernameTaken(sqlDB, username)
	if reqErr != nil {
		return reqErr
	}

	if (taken) {
		return user_errors.UsernameTakenErr(errors.New("username already taken"))
	}
	
	err := u.SetUsername(uId, username)
	
	if err != nil {
		return user_errors.SetNameFailedErr(err)
	}

	return nil
}

func setPageTheme(
	sqlDB *sql.DB, 
	uId int,
	pageTheme string, 
) (*models.RequestError) {

	u := user_db.UserDB{DB: sqlDB}
	
	err := u.SetPageTheme(uId, pageTheme)
	
	if err != nil {
		return user_errors.SetPageThemeFailedErr(err)
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

	if (user.BankConnected) {
		return banking_errors.AlreadyInitialisedBankingErr(err)
	}

	// 1. Get access token
	accessToken, err := my_plaid.GetAuthAccessToken(plaidCli, publicToken)
	// fmt.Println("Access token :", accessToken)
	if err != nil {
		return banking_errors.GetAccessTokenFailedErr(err)
	}

	err = u.SetAccessToken(uId, accessToken)
	if err != nil {
		return user_errors.SetAccessTokenTokenFailedErr(err)
	}

	// 2. Get account bank details
	bacs, err := my_plaid.GetBACNumbers(plaidCli, accessToken)
	// fmt.Println("Bacs sort code:", bacs.SortCode)
	if err != nil {
		return banking_errors.GetBACSFailedErr(err)
	}
	
	recipientID, err := my_plaid.CreatePaymentRecipient(plaidCli, user.FullName(), bacs.Account, bacs.SortCode)
	// fmt.Println("Recipient ID:", recipientID)
	if err != nil {
		return banking_errors.CreatePaymentRecipientFailedErr(err)
	}

	// 1. Set public token
	err = u.SetRecipientID(uId, recipientID)
	// fmt.Println("Successfuly set public token")
	if err != nil {
		return user_errors.SetRecipientIDFailedErr(err)
	}

	// Call transactions/sync 
	go my_plaid.SyncTransactions(plaidCli, accessToken, nil)

	return nil
}