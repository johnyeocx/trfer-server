package models

type EmailOTP struct {
	Email			string 	`json:"email"`
	HashedOTP		string 	`json:"hashed_otp"`
}