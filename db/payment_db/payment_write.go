package payment_db

import (
	"database/sql"

	"github.com/johnyeocx/usual/server2/db/models"
)

type PaymentDB struct {
	DB	*sql.DB
}


func (p *PaymentDB) InsertPayment(
	sqlDB *sql.DB, 
	payment models.Payment,
) (error) {
	query := `INSERT INTO payment 
	(plaid_payment_id, username, amount, note, created, payment_status) VALUES 
	($1, $2, $3, $4, $5, $6)`

	_, err := p.DB.Query(query, 
		payment.PlaidPaymentID, 
		payment.Username,
		payment.Amount, 
		payment.Note, 
		payment.Created,
		payment.PaymentStatus,
	)

	return err
}