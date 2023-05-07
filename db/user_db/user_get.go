package user_db

import (
	"database/sql"

	"github.com/johnyeocx/usual/server2/db/models/user_models"
)


type UserDB struct {
	DB	*sql.DB
}


func ValidateUser (sqlDB *sql.DB, id int) (bool) {
	var email string

	err := sqlDB.QueryRow(`SELECT email FROM "user" WHERE user_id=$1 AND email_verified='True'`, 
		id,
	).Scan(&email) 
	
	if err != nil {
		return false
	}
	
	return email != ""
}


func (u *UserDB) GetUserByUsername(username string) (*user_models.User, error) {
	var user user_models.User

	err := u.DB.QueryRow(`SELECT 
		user_id, username, email, first_name, last_name, public_token, recipient_id, page_theme
		FROM "user" WHERE username=$1 AND email_verified=TRUE`, 
	username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PublicToken,
		&user.RecipientID,
		&user.PageTheme,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDB) GetUserByEmail(email string) (*user_models.User, error) {
	var user user_models.User
	
	err := u.DB.QueryRow(`SELECT 
		user_id, username, email, first_name, last_name
		FROM "user" WHERE email=$1 AND email_verified=TRUE`, 
	email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDB) GetUserByID(uId int) (*user_models.User, error) {
	var user user_models.User
	
	err := u.DB.QueryRow(`SELECT 
		user_id, username, email, first_name, last_name, public_token, page_theme
		FROM "user" WHERE user_id=$1 AND email_verified=TRUE`, 
	uId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PublicToken,
		&user.PageTheme,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDB) GetLastPaymentID(uId int) (int, error) {
	lastPaymentId := 0
	err := u.DB.QueryRow(`SELECT 
		payment_id FROM payment WHERE user_id=$1 ORDER BY payment_id DESC LIMIT 1`, 
	uId).Scan(&lastPaymentId)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return -1, err
	}

	return lastPaymentId, nil
}