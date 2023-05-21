package pers_db

import (
	"database/sql"

	personamodels "github.com/johnyeocx/usual/server2/db/models/persona_models"
)

type PersDB struct {
	DB	*sql.DB
}

func (i *PersDB) GetUserLatestInquiry(uId int) (*personamodels.Inquiry, error) {
	
	query := `
		SELECT pers_inquiry_id, i.pers_account_id, pers_session_id, inquiry_status, session_status
		from "user" as u JOIN inquiry as i
		on u.pers_account_id=i.pers_account_id
		WHERE u.user_id=$1
		ORDER BY created_at DESC LIMIT 1
	`
	
	inq := personamodels.Inquiry{}
	err := i.DB.QueryRow(query, uId).Scan(
		&inq.PersInquiryID,
		&inq.PersAccountID,
		&inq.PersSessionID,
		&inq.InquiryStatus,
		&inq.SessionStatus,
	)
	if err != nil {
		return nil, err
	}

	return &inq, err
}

func (i *PersDB) InsertInquiry(inquiry personamodels.Inquiry) (error) {
	
	startedAt := sql.NullTime{}
	if (inquiry.StartedAt != nil) {
		startedAt.Valid = true
		startedAt.Time = *inquiry.StartedAt
	}

	iSqlNulls := inquiry.GetSqlNullVals()
	
	_, err := i.DB.Exec(`INSERT into inquiry 
		(pers_inquiry_id, pers_account_id, inquiry_status, 
			created_at, started_at, completed_at, decisioned_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT (pers_inquiry_id) DO UPDATE 
		SET inquiry_status=$3
		` , 
		inquiry.PersInquiryID, 
		inquiry.PersAccountID, 
		inquiry.InquiryStatus, 

		iSqlNulls.CreatedAt,
		iSqlNulls.StartedAt,
		iSqlNulls.CompletedAt,
		iSqlNulls.DecisionedAt,
	)

	return err
}

func (i *PersDB) UpdateInquirySession(inquiry personamodels.Inquiry) (error) {
	
	_, err := i.DB.Exec(`UPDATE inquiry SET pers_session_id=$1, session_status=$2
		WHERE pers_inquiry_id=$3
		` , 
		inquiry.PersSessionID, 
		inquiry.SessionStatus,
		inquiry.PersInquiryID, 
	)

	return err
}