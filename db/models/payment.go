package models

import (
	"time"

	"github.com/johnyeocx/usual/server2/utils/enums"
)

type Payment struct {
	ID 					int 					`json:"payment_id"`
	PlaidPaymentID 		string 					`json:"plaid_payment_id"`
	UserID		 		int 					`json:"user_id"`
	Amount				float64					`json:"amount"`
	Note				string				 	`json:"note"`
	Reference			string 					`json:"reference"`
	Created 			time.Time				`json:"created"`
	PaymentStatus		enums.PaymentStatus		`json:"payment_status"`
	AccessToken			JsonNullString			`json:"access_token"`
	TransactionName		JsonNullString			`json:"transaction_name"`
}

type PaymentInitiationEvent struct {
	NewPaymentStatus	enums.PaymentStatus `json:"new_payment_status`
	PlaidPaymentID 		string 				`json:"payment_id"`
}