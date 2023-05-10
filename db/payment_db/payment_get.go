package payment_db

import (
	"github.com/johnyeocx/usual/server2/db/models"
)

func (p *PaymentDB) GetPayment(plaidPaymentId string) (*models.Payment, error) {
	query := `SELECT payment_id, reference, user_id, created FROM payment WHERE plaid_payment_id=$1`
	
	payment := models.Payment{}
	err := p.DB.QueryRow(query, plaidPaymentId).Scan(
		&payment.ID,
		&payment.Reference,
		&payment.UserID,
		&payment.Created,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (p *PaymentDB) GetUnnamedPayments() ([]models.Payment, error) {
	query := `
		SELECT 
		p.payment_id, p.reference, p.user_id, p.created, u.access_token FROM 
		payment as p JOIN "user" as u on p.user_id=u.user_id
		WHERE "transaction_name" is NULL
		ORDER BY p.user_id, p.created ASC
	`
	
	
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}


	defer rows.Close()

	results := []models.Payment{}

	for rows.Next() {
		payment := models.Payment{}

        if err := rows.Scan(
			&payment.ID,
			&payment.Reference,
			&payment.UserID,
			&payment.Created,
			&payment.AccessToken,
		); err != nil {
			continue
        }

		results = append(results, payment)
    }

	return results, nil
}