package user_db

import (
	"database/sql"

	"github.com/johnyeocx/usual/server2/db/models/user_models"
)

func (u *UserDB) CreateUserFromEmail (
	email string,
	username string,
	verified bool,
	persActId *string,
) (*int, error) {

	persId := sql.NullString{}
	if persActId != nil {
		persId.Valid = true
		persId.String = *persActId
	}

	_, err := u.DB.Exec(`DELETE from "user" WHERE email=$1 AND email_verified=FALSE`, email)
	if err != nil {
		return nil, err
	}

	_, err = u.DB.Exec(`DELETE from "user" WHERE username=$1 AND email_verified=FALSE`, username)
	if err != nil {
		return nil, err
	}

	var userId int
	err = u.DB.QueryRow(`
		INSERT into "user" (email, username, email_verified, pers_account_id) VALUES ($1, $2, $3, $4)
		RETURNING user_id`,
		email, username, verified, persActId,
	).Scan(&userId)

	if err != nil {
		return nil, err
	}

	return &userId, nil
}

func (u *UserDB) SetPersIdAndEmailverified(persId string, email string) (*user_models.User, error) {
	var user user_models.User
	
	err := u.DB.QueryRow(`UPDATE "user" SET pers_account_id=$1, email_verified=TRUE WHERE email=$2 RETURNING user_id`, 		
		persId, 
		email,
	).Scan(
		&user.ID,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDB) SetName(id int, firstName string, lastName string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET first_name=$1, last_name=$2 WHERE user_id=$3`,
		firstName, lastName, id,
	)

	return err
}

func (u *UserDB) SetUsername(id int, username string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET username=$1 WHERE user_id=$2`,
		username, id,
	)

	return err
}


func (u *UserDB) SetPageTheme(id int, pageTheme string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET page_theme=$1 WHERE user_id=$2`,
		pageTheme, id,
	)

	return err
}

func (u *UserDB) SetRecipientID(id int, recipientId string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET recipient_id=$1 WHERE user_id=$2`,
		recipientId, id,
	)

	return err
}

func (u *UserDB) SetAccessToken(id int, accessToken string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET access_token=$1, bank_connected='True' WHERE user_id=$2`,
		accessToken, id,
	)

	return err
}