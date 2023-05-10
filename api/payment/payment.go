package payment

import (
	"database/sql"
	"fmt"
	"log"
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
	reference := fmt.Sprintf("TMU%dP%d", user.ID, lastPaymentId)
	
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
	sqlDB *sql.DB, 
	plaidCli *plaid.APIClient, 
	piEvent models.PaymentInitiationEvent,
) (*models.RequestError) {
	p := payment_db.PaymentDB{DB: sqlDB}

	payerName := models.JsonNullString{
		sql.NullString{
			Valid: false,
			String: "",
		},
	}

	if (piEvent.NewPaymentStatus == PaymentStatus.Executed) {
		payment, err := my_plaid.GetPayment(plaidCli, piEvent.PlaidPaymentID)
		if err != nil {
			return banking_errors.GetPaymentFailedErr(err)
		}
	
		refundDetails := payment.RefundDetails
	
	
		if (refundDetails.IsSet() && refundDetails.Get().Name != "") {
			payerName.Valid = true
			payerName.String = refundDetails.Get().GetName()
		}	
	}

	err := p.UpdatePaymentFromPIEvent(piEvent, payerName)
	if err != nil {
		return banking_errors.UpdatePaymentFailedErr(err)
	}

	return nil
}

func UpdatePaymentNames(
	sqlDB *sql.DB, 
	plaidCli *plaid.APIClient,
	userId int,
	accessToken string,
	startDate time.Time,
	endDate time.Time,
	payments []models.Payment,
) (error) {
	txs, err := my_plaid.GetUserTransactions(plaidCli, accessToken, startDate, endDate)
	if err != nil {
		return err
	}
	

	refToName := map[string]string{}
	for _, tx := range(txs) {
		if tx.PaymentMeta.ReferenceNumber.IsSet() {
			refToName[*tx.PaymentMeta.ReferenceNumber.Get()] = tx.Name
		}
		fmt.Printf("Name: %s, Reference: %v\n", tx.Name, tx.PaymentMeta.ReferenceNumber)
	}
	fmt.Println(refToName)

	namedPayments := []models.Payment{}
	for _, payment := range(payments) {
		ref := payment.Reference
		fmt.Println(ref)
		if val, ok := refToName[ref]; ok {
			payment.TransactionName = models.JsonNullString{
				sql.NullString{
					Valid: true,
					String: val,
				},
			}
			namedPayments = append(namedPayments, payment)
		}
	}

	if len(namedPayments) > 0 {
		pDB := payment_db.PaymentDB{DB: sqlDB}
		err := pDB.UpdatePaymentNames(namedPayments)
		if err != nil {
			fmt.Println("Failed to update payment named:", err)
		}
	} else {
		log.Println("No payments to update")
	}

	return nil
}