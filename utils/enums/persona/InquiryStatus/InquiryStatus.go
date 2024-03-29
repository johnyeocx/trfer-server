package InquiryStatus

import (
	"github.com/johnyeocx/usual/server2/utils/enums"
)


const (
	Created enums.InquiryStatus = "created"
	Pending enums.InquiryStatus = "pending"
	Failed enums.InquiryStatus = "failed"
	NeedsReview enums.InquiryStatus = "needs_review"

	Completed enums.InquiryStatus = "completed"
	Approved enums.InquiryStatus = "approved"
	Declined enums.InquiryStatus = "declined"
)

// func EventStrToPaymentStatus(str string) (enums.PaymentStatus) {
// 	if (str == "PAYMENT_STATUS_AUTHORISING") {
// 		return Authorising
// 	} else if (str == "PAYMENT_STATUS_EXECUTED") {
// 		return Executed
// 	} else if (str == "PAYMENT_STATUS_INPUT_NEEDED") {
// 		return InputNeeded
// 	} else if (str == "PAYMENT_STATUS_REJECTED") {
// 		return Rejected
// 	}
// 	return Created

// 	// PAYMENT_STATUS_INPUT_NEEDED
// }