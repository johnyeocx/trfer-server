package banking

import (
	"database/sql"
	"errors"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)


func setUserAccessToken(
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

	if (user.AccessToken.Valid) {
		return banking_errors.AlreadySetAccessTokenErr(err)
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

	return nil
}

func createRecipientID(sqlDB *sql.DB, plaidCli *plaid.APIClient, uId int) (*models.RequestError){
	u := user_db.UserDB{DB: sqlDB}

	user, err := u.GetUserByID(uId)
	if err != nil {
		return user_errors.GetUserFailedErr(err)
	}
	if (!user.AccessToken.Valid) {
		return banking_errors.NoAccessTokenErr(errors.New("no access token present"))
	}

	// 2. Get account bank details
	bacs, err := my_plaid.GetBACNumbers(plaidCli, user.AccessToken.String)

	// fmt.Println("Bacs sort code:", bacs.SortCode)
	if err != nil {
		return banking_errors.GetBACSFailedErr(err)
	}
	
	recipientID, err := my_plaid.CreatePaymentRecipient(plaidCli, user.FullName(), bacs.Account, bacs.SortCode)
	if err != nil {
		return banking_errors.CreatePaymentRecipientFailedErr(err)
	}

	// 1. Set public token
	err = u.SetRecipientID(uId, recipientID)
	if err != nil {
		return user_errors.SetRecipientIDFailedErr(err)
	}

	return nil
}