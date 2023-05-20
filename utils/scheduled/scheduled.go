package scheduled

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/johnyeocx/usual/server2/db/payment_db"
	"github.com/plaid/plaid-go/v11/plaid"
)

func RunScheduled(sqlDB *sql.DB, plaidCli *plaid.APIClient) {

	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minutes().Do(func() {
		PollForTransactions(sqlDB, plaidCli)
	})
	s.StartBlocking()
}

func PollForTransactions(sqlDB *sql.DB, plaidCli *plaid.APIClient) {

	p := payment_db.PaymentDB{DB: sqlDB}
	payments, err := p.GetUnnamedPayments()
	if err != nil {
		return
	}


	startDate := time.Now()
	endDate := time.Now()
	startIndex := 0
	endIndex := 0

	for i := 0; i < len(payments); i++ {
		
		if i == 0 || payments[i].UserID != payments[i - 1].UserID {
			// is first for user
			created := payments[i].Created
			startDate = time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, time.UTC)
			startIndex = i
		}

		if i == len(payments) - 1 || payments[i].UserID != payments[i + 1].UserID {
			created := payments[i].Created
			endDate = time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
			endIndex = i
			
			// Exec search function on start and end
			log.Printf("Updating transaction names for user %d, from %v to %v\n", payments[i].UserID, startDate.Format("01/02"), endDate.Format("01/02"))
			fmt.Println(startIndex, endIndex)
			// err := payment.UpdatePaymentNames(
			// 	sqlDB, 
			// 	plaidCli, 
			// 	payments[i].UserID, 
			// 	payments[i].AccessToken.String, 
			// 	startDate, 
			// 	endDate,
			// 	payments[startIndex: endIndex],
			// )
			// if err != nil {
			// 	log.Println("Failed to update payment names:", err)
			// }
			// fmt.Println()
		}		
	}

	return
}

