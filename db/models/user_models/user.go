package user_models

import (
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/utils/enums"
)

type User struct {
	ID 					int 					`json:"customer_id"`
	Email 				string 					`json:"email"`
	Username 			string 					`json:"username"`
	FirstName			models.JsonNullString	`json:"first_name"`
	LastName			models.JsonNullString 	`json:"last_name"`
	BankConnected		bool					`json:"bank_connected"`
	AccessToken			models.JsonNullString 	`json:"access_token"`
	RecipientID			models.JsonNullString 	`json:"recipient_id"`
	SignInProvider		enums.SignInProvider	`json:"signin_provider"`
	PageTheme 			string			`json:"page_theme"`
}

func (u* User) FullName() string{
	return u.FirstName.String + " " + u.LastName.String
}

