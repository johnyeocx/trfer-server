package otp_db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/utils/enums"
)


type OtpDB struct {
	DB	*sql.DB
}

func (o *OtpDB) GetEmailOtp(email string, otpType enums.OtpType) (*models.EmailOTP, error) {

	fmt.Println("OTP TYPE:", otpType)
	selectStatement := `
		SELECT hashed_otp, email from email_otp WHERE
		email=$1 AND type=$2 AND $3 <= expiry
	`

	row := o.DB.QueryRow(selectStatement, email, string(otpType), time.Now().UTC())

	var emailOtp models.EmailOTP
	err := row.Scan(&emailOtp.HashedOTP, &emailOtp.Email);

	if err == sql.ErrNoRows {
		return nil, err
	}

	return &emailOtp, nil
}

func (o *OtpDB) DeleteEmailOtp(email string, otpType enums.OtpType) (error) {

    deleteStatement := `
        DELETE from email_otp WHERE email=$1 AND type=$2
    `
	
    _, err := o.DB.Exec(deleteStatement, email, string(otpType))

	if err != nil {
		return err
	}

	return nil
}

func (o *OtpDB) InsertEmailOtp(email string, otpType enums.OtpType, hashedOtp string, expiry time.Time) (error) {

	_, err := o.DB.Exec(
		`
			INSERT INTO email_otp (email, type, hashed_otp, expiry) 
			VALUES ($1, $2, $3, $4)
		`,
		email, string(otpType), hashedOtp, expiry,
	)

	return err
}