package PaymentStatus

import (
	"github.com/johnyeocx/usual/server2/utils/enums"
)


const (
	Created enums.PaymentStatus = "created"
	Authorising enums.PaymentStatus = "authorising"
	InputNeeded enums.PaymentStatus = "input_needed"
	Executed enums.PaymentStatus = "executed"
	Rejected enums.PaymentStatus = "rejected"
	
)

func EventStrToPaymentStatus(str string) (enums.PaymentStatus) {
	if (str == "PAYMENT_STATUS_AUTHORISING") {
		return Authorising
	} else if (str == "PAYMENT_STATUS_EXECUTED") {
		return Executed
	} else if (str == "PAYMENT_STATUS_INPUT_NEEDED") {
		return InputNeeded
	} else if (str == "PAYMENT_STATUS_REJECTED") {
		return Rejected
	}
	return Created

	// PAYMENT_STATUS_INPUT_NEEDED
}