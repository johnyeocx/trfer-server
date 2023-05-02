package transfer

import (
	"database/sql"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/payment_db"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/banking_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	"github.com/johnyeocx/usual/server2/utils/enums/PaymentStatus"
	my_plaid "github.com/johnyeocx/usual/server2/utils/plaid"
	"github.com/plaid/plaid-go/v11/plaid"
)

func TransferOpenAmt(sqlDB *sql.DB, plaidCli *plaid.APIClient, username string, amount float64, note string) (string, *models.RequestError) {

	// 1. Get recipient id
	u := user_db.UserDB{DB: sqlDB}
	user, err := u.GetUserByUsername(username)
	if err != nil {
		return "", user_errors.GetUserFailedErr(err)
	}

	// 2. Create payment request
	paymentID, err := my_plaid.CreatePayment(plaidCli, user.RecipientID.String, amount)
	if err != nil {
		return "", banking_errors.CreatePaymentFailedErr(err)
	}

	// 3. Create link token
	linkToken, err := my_plaid.CreatePaymentLinkToken(plaidCli, user.ID, paymentID)
	if err != nil {
		return "", banking_errors.CreateAuthLinkTokenFailedErr(err)
	}

	// 4. Add payment to DB
	p := payment_db.PaymentDB{DB: sqlDB}
	err = p.InsertPayment(sqlDB, models.Payment{
		PlaidPaymentID: paymentID,
		Username: user.Username,
		Amount: amount,
		Note: note,
		Created: time.Now(),
		PaymentStatus: PaymentStatus.Created,
	})
	
	if err != nil {
		return "", banking_errors.InsertPaymentFailedErr(err)
	}

	return linkToken, nil
}