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

func (u *UserDB) SetAccountName(id int, accountName string) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET account_name=$1 WHERE user_id=$2`,
		accountName, id,
	)

	return err
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
		`UPDATE "user" SET access_token=$1 WHERE user_id=$2`,
		accessToken, id,
	)

	return err
}

func (u *UserDB) SetPersApproved(acctId string, persApproved bool) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET pers_approved=$1 WHERE pers_account_id=$2`,
		persApproved, acctId,
	)

	return err
}

func (u *UserDB) SetAddress(address user_models.Address, uId int) (error) {
	
	_, err := u.DB.Exec(
		`UPDATE "user" SET line1=$1, line2=$2, city=$3, postal_code=$4 WHERE user_id=$5`,
		address.Line1,
		address.Line2,
		address.City,
		address.PostalCode,
		uId,
	)

	return err
}