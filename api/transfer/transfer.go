package transfer

import (
	"database/sql"
	"fmt"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)

func TransferOpenAmt(sqlDB *sql.DB, plaidCli *plaid.APIClient, username string, amount int) (string, *models.RequestError) {
	// 1. Get recipient id
	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByUsername(username)
	if err != nil {
		return "", user_errors.GetUserFailedErr(err)
	}

	// 2. Create payment request
	amountFloat := float64(amount)
	amountFloat = amountFloat / 2
	fmt.Println("Amount: ", amountFloat)

	paymentID, err := my_plaid.CreatePayment(plaidCli, user.RecipientID.String, amountFloat)
	if err != nil {
		return "", banking_errors.CreatePaymentFailedErr(err)
	}

	// 3. Create link token
	linkToken, err := my_plaid.CreatePaymentLinkToken(plaidCli, user.ID, paymentID)
	if err != nil {
		return "", banking_errors.CreateAuthLinkTokenFailedErr(err)
	}

	return linkToken, nil
}