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
	PublicToken			models.JsonNullString 	`json:"public_token"`
	RecipientID		models.JsonNullString 	`json:"recipient_id"`
	SignInProvider		enums.SignInProvider	`json:"signin_provider"`
}

func (u* User) FullName() string{
	return u.FirstName.String + " " + u.LastName.String
}

