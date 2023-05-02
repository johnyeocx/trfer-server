package payment

import (
	"database/sql"
	"fmt"
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

	lastPaymentId, err := u.GetLastPaymentID(user.ID)
	if err != nil {
		return "", user_errors.GetLastPaymentIDFailedErr(err)
	}

	
	// create reference (could cause potential error in race condition (?))
	reference := fmt.Sprintf("tmu%dp%d", user.ID, lastPaymentId)
	
	// 2. Create payment request
	paymentID, err := my_plaid.CreatePayment(plaidCli, user.RecipientID.String, amount, reference)
	if err != nil {
		fmt.Println("Failed to create payment:", err)
		return "", banking_errors.CreatePaymentFailedErr(err)
	}

	// 3. Create link token
	linkToken, err := my_plaid.CreatePaymentLinkToken(plaidCli, user.ID, paymentID)
	if err != nil {
		return "", banking_errors.CreateAuthLinkTokenFailedErr(err)
	}

	// 4. Add payment to DB
	p := payment_db.PaymentDB{DB: sqlDB}
	err = p.InsertPayment(models.Payment{
		PlaidPaymentID: paymentID,
		UserID: user.ID,
		Amount: amount,
		Reference: reference,
		Note: note,
		Created: time.Now(),
		PaymentStatus: PaymentStatus.Created,
	})
	
	if err != nil {
		return "", banking_errors.InsertPaymentFailedErr(err)
	}

	return linkToken, nil
}

func UpdatePaymentFromPIEvent(
	sqlDB *sql.DB, plaidCli *plaid.APIClient, piEvent models.PaymentInitiationEvent) (*models.RequestError) {
	p := payment_db.PaymentDB{DB: sqlDB}
	u := user_db.UserDB{DB: sqlDB}

	txName := sql.NullString{
		Valid: false,
	}

	if piEvent.NewPaymentStatus == PaymentStatus.Executed {
		payment, err := p.GetPayment(piEvent.PlaidPaymentID)
	if err != nil {
		return banking_errors.GetPaymentFailedErr(err)
	}

	user, err := u.GetUserByID(payment.UserID)
	if err != nil {
		return user_errors.GetUserFailedErr(err)
	}

		transactions, err := my_plaid.GetUserTransactions(plaidCli, user.PublicToken.String, payment.Created)
		if err != nil {
			return banking_errors.GetTransactionsFailedErr(err)
		}

		for _, transaction := range(transactions) {
			txReferenceNumber := transaction.PaymentMeta.ReferenceNumber
			if *txReferenceNumber.Get() == payment.Reference {
				txName.Valid = true
				txName.String = transaction.Name
			}
		}
	}

	if (txName.Valid) {
		fmt.Println("Transaction Name:", txName.String)
	} else {
		fmt.Println("No transaction name")
	}
	
	err := p.UpdatePaymentFromPIEvent(piEvent, txName)
	if err != nil {
		return banking_errors.UpdatePaymentFailedErr(err)
	}

	return nil
}