package user_models

import "github.com/johnyeocx/usual/server2/db/models"

type Address struct {
	Line1 			models.JsonNullString `json:"line1"`
	Line2 			models.JsonNullString `json:"line2"`
	City 			models.JsonNullString `json:"city"`
	PostalCode 		models.JsonNullString `json:"postal_code"`
	CountryCode 	models.JsonNullString `json:"country_code"`
}