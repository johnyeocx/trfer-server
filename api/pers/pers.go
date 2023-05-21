package persona

import (
	"database/sql"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	personamodels "github.com/johnyeocx/usual/server2/db/models/persona_models"
	"github.com/johnyeocx/usual/server2/db/pers_db"
	"github.com/johnyeocx/usual/server2/errors/pers_errors"
)


func DecodeInquiryWebhook(data map[string]interface{}) (*personamodels.Inquiry, *models.RequestError){
	inqId := data["id"].(string)

	attributes := data["attributes"].(map[string]interface{})
	status := attributes["status"].(string)

	createdAtStr := attributes["created-at"].(string)
	_, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, pers_errors.DecodeInquiryFailedErr(err)
	}

	relationships := data["relationships"].(map[string]interface{})
	account := relationships["account"].(map[string]interface{})["data"].(map[string]interface{})
	acctId := account["id"].(string)
	
	return &personamodels.Inquiry{
		PersInquiryID: inqId,
		InquiryStatus: status,
		PersAccountID: acctId,
	}, nil
}

func UpdateInquiry(sqlDB *sql.DB, inquiry personamodels.Inquiry) {
	p := pers_db.PersDB{DB: sqlDB}
	p.InsertInquiry(inquiry)
}