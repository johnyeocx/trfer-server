package cus_models

import (
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/utils/enums"
)

type Customer struct {
	ID 					int 					`json:"customer_id"`
	Email 				string 					`json:"email"`
	Password 			models.JsonNullString	`json:"password"`
	EmailVerified		bool 					`json:"email_verified"`
	FirstName			models.JsonNullString	`json:"first_name"`
	LastName			models.JsonNullString 	`json:"last_name"`
	SignInProvider		enums.SignInProvider	`json:"signin_provider"`
	// Address 			*CusAddress 	`json:"address"`
	// StripeID 			string 			`json:"stripe_id"`
	// DefaultCardID	 	JsonNullInt16	`json:"default_card_id"`
	// Uuid 				string 			`json:"uuid"`
}

// func (c* Customer) FullName() string{
// 	return constants.FullName(c.FirstName, c.LastName)
// }

