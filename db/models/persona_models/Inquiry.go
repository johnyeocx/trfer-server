package personamodels

import (
	"database/sql"
	"time"
)

type Inquiry struct {
	ID string `json:"id"`
	PersInquiryID string `json:"pers_inquiry_id"`
	PersAccountID string `json:"pers_account_id"`
	InquiryStatus string `json:"inquiry_status"`
	PersSessionID string `json:"pers_session_id"`
	SessionStatus string `json:"session_status"`

	CreatedAt 			*time.Time `json:"created_at"`
	StartedAt 			*time.Time `json:"started_at"`
	CompletedAt 		*time.Time `json:"completed_at"`
	DecisionedAt 		*time.Time `json:"decisioned_at"`	
	// FailedAt 			*time.Time `json:"failed_at"`
	// MarkedForReviewAt 	*time.Time `json:"marked_for_review_at"`	
	// ExpiredAt 			*time.Time `json:"expired_at"`	
	// RedactedAt 			*time.Time `json:"redacted_at"`		
}

type InquirySqlNulls struct {
	CreatedAt 		sql.NullTime `json:"created_at"`
	StartedAt 		sql.NullTime `json:"started_at"`
	CompletedAt 	sql.NullTime `json:"completed_at"`
	DecisionedAt 	sql.NullTime `json:"decisioned_at"`	
}
func (i *Inquiry) GetSqlNullVals() (InquirySqlNulls){
	iSql := InquirySqlNulls{}
	if (i.CreatedAt != nil) {
		iSql.CreatedAt.Valid = true
		iSql.CreatedAt.Time = *i.CreatedAt
	}

	if (i.StartedAt != nil) {
		iSql.StartedAt.Valid = true
		iSql.StartedAt.Time = *i.StartedAt
	}
	if (i.CompletedAt != nil) {
		iSql.CompletedAt.Valid = true
		iSql.CompletedAt.Time = *i.CompletedAt
	}
	if (i.DecisionedAt != nil) {
		iSql.DecisionedAt.Valid = true
		iSql.DecisionedAt.Time = *i.DecisionedAt
	}
	return iSql
}