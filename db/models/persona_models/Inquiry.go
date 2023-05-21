package personamodels

type Inquiry struct {
	ID string `json:"id"`
	PersInquiryID string `json:"pers_inquiry_id"`
	PersAccountID string `json:"pers_account_id"`
	InquiryStatus string `json:"inquiry_status"`
	PersSessionID string `json:"pers_session_id"`
	SessionStatus string `json:"session_status"`
}