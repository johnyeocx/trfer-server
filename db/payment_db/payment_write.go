package payment_db

import (
	"database/sql"
	"fmt"

	"github.com/johnyeocx/usual/server2/db/models"
)

type PaymentDB struct {
	DB	*sql.DB
}


func (p *PaymentDB) InsertPayment(
	payment models.Payment,
) (error) {
	query := `INSERT INTO payment 
	(plaid_payment_id, user_id, amount, note, created, payment_status, reference) VALUES 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.DB.Query(query, 
		payment.PlaidPaymentID, 
		payment.UserID,
		payment.Amount, 
		payment.Note, 
		payment.Created,
		payment.PaymentStatus,
		payment.Reference,
	)

	return err
}


func (p *PaymentDB) UpdatePaymentFromPIEvent(
	piEvent models.PaymentInitiationEvent,
	name models.JsonNullString,
) (error) {
	query := `UPDATE payment SET payment_status=$1, payer_name=$2 WHERE plaid_payment_id=$3`

	_, err := p.DB.Exec(query, 
		piEvent.NewPaymentStatus, 
		name,
		piEvent.PlaidPaymentID,
	)

	return err
}

func (p *PaymentDB) UpdatePaymentNames(
	// piEvent models.PaymentInitiationEvent,
	payments []models.Payment,
) (error) {

	valuesString := ``
	for i, payment := range(payments) {
		valuesString += fmt.Sprintf(`(%d, '%s')`, payment.ID, payment.TransactionName.String)
		if i != len(payments) - 1 {
			valuesString += ", "
		}
	}

	
	queryString := fmt.Sprintf(`
	UPDATE payment as p SET transaction_name=c.transaction_name FROM
	(VALUES 
		%s 
	)
	AS c(payment_id, transaction_name)
	WHERE c.payment_id=p.payment_id
	`, valuesString)

	fmt.Println(queryString)

	p.DB.Exec(queryString)

	return nil
}