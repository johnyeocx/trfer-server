package models

import (
	"time"

	"github.com/johnyeocx/usual/server2/utils/enums"
)

type Payment struct {
	ID 					int 					`json:"payment_id"`
	PlaidPaymentID 		string 					`json:"plaid_payment_id"`
	Username		 	string 					`json:"username"`
	Amount				float64					`json:"amount"`
	Note				string				 	`json:"note"`
	Created 			time.Time				`json:"created"`
	PaymentStatus		enums.PaymentStatus		`json:"payment_status"`
}