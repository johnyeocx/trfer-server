package personamodels

type Account struct {
	ID string `json:"id"`
	Attributes map[string]interface{}
}

type Attributes struct {
	// ReferenceID string 	`json:"reference_id"`
	// CreatedAt time.Time `json:"created_at"`
	NameFirst string 	`json:"name-first"`
	NameMiddle string 	`json:"name-middle"`
}