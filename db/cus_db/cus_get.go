package cus_db

import (
	"database/sql"

	"github.com/johnyeocx/usual/server2/db/models/cus_models"
)


type CustomerDB struct {
	DB	*sql.DB
}


func (c *CustomerDB) GetCustomerByEmail (email string) (*cus_models.Customer, error) {
	var cus cus_models.Customer 
	err := c.DB.QueryRow(`SELECT 
		customer_id, email, email_verified, signin_provider, first_name, last_name
		FROM customer WHERE email=$1`, 
	email).Scan(
		&cus.ID,
		&cus.Email,
		&cus.EmailVerified,
		&cus.SignInProvider,
		&cus.FirstName,
		&cus.LastName,
	)

	if err != nil {
		return nil, err
	}

	return &cus, nil
}

func ValidateCustomerId (sqlDB *sql.DB, id int) (bool) {
	var email string

	err := sqlDB.QueryRow("SELECT email FROM customer WHERE customer_id=$1", 
		id,
	).Scan(&email) 
	
	if err != nil {
		return false
	}
	

	return email != ""
}

