package sub_models

type Subscription struct {
	ID 					int				`json:"sub_id"`
	Brand				Brand			`json:"brand"`
	Plan				Plan			`json:"plan"`
	Status				string			`json:"status"`
}