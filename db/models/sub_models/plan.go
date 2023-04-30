package sub_models

import (
	"time"

	"github.com/johnyeocx/usual/server2/utils/enums"
)

type Plan struct {
	Interval		enums.Interval		`json:"interval"`
	UnitAmount		int 				`json:"unit_amount"`
	IntervalCount	int 				`json:"interval_count"`
	Currency		string 				`json:"currency"`
	BillDate		time.Time			`json:"bill_date"`
}