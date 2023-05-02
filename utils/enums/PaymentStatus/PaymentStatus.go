package PaymentStatus

import "github.com/johnyeocx/usual/server2/utils/enums"


const (
	Created enums.PaymentStatus = "created"
	Authorising enums.PaymentStatus = "authorising"
	Executed enums.PaymentStatus = "executed"
)

func EventStrToPaymentStatus(str string) (enums.PaymentStatus) {
	if (str == "PAYMENT_STATUS_AUTHORISING") {
		return Authorising
	} else if (str == "PAYMENT_STATUS_EXECUTED") {
		return Executed
	} else {
		return Created
	}
}