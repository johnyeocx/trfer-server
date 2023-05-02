package payment_db

import "github.com/johnyeocx/usual/server2/db/models"

func (p *PaymentDB) GetPayment(plaidPaymentId string) (*models.Payment, error) {
	query := `SELECT payment_id, reference, user_id, created FROM payment WHERE plaid_payment_id=$1`
	
	payment := models.Payment{}
	err := p.DB.QueryRow(query, plaidPaymentId).Scan(
		payment.ID,
		payment.Reference,
		payment.UserID,
		payment.Created,
	)

	if err != nil {
		return nil, nil
	}

	return &payment, nil
}