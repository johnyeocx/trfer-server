package pers_db

import (
	"database/sql"

	personamodels "github.com/johnyeocx/usual/server2/db/models/persona_models"
)

type PersDB struct {
	DB	*sql.DB
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
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (pers_inquiry_id) DO UPDATE 
		SET inquiry_status=$3
		` , 
		inquiry.PersInquiryID, 
		inquiry.PersAccountID, 
		inquiry.PersSessionID, 
		inquiry.SessionStatus,

		iSqlNulls.CreatedAt,
		iSqlNulls.StartedAt,
		iSqlNulls.CompletedAt,
		iSqlNulls.DecisionedAt,
	)

	return err
}