package cus_db

import "github.com/johnyeocx/usual/server2/utils/enums"

func (c *CustomerDB) CreateCustomerFromExtSignin (
	email string,
	signinProvider enums.SignInProvider,
) (*int, error) {

	var cusId int
	err := c.DB.QueryRow(`
		INSERT into customer (email, signin_provider, email_verified) VALUES ($1, $2, 'True') 
		ON CONFLICT (email)  DO UPDATE 
		SET signin_provider=$2, email_verified='True', first_name='', last_name=''
		RETURNING customer_id`,
		email, signinProvider,
	).Scan(&cusId)

	if err != nil {
		return nil, err
	}

	return &cusId, nil
}
