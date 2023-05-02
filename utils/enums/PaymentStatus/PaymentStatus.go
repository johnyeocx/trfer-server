package PaymentStatus

import "github.com/johnyeocx/usual/server2/utils/enums"


const (
	Created enums.PaymentStatus = "created"
	Authorising enums.PaymentStatus = "authorising"
	Executed enums.PaymentStatus = "executed"
)