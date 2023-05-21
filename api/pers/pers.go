package persona

import (
	"database/sql"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	personamodels "github.com/johnyeocx/usual/server2/db/models/persona_models"
	"github.com/johnyeocx/usual/server2/db/pers_db"
)


func DecodeInquiryWebhook(data map[string]interface{}) (*personamodels.Inquiry, *models.RequestError){
	inqId := data["id"].(string)

	attributes := data["attributes"].(map[string]interface{})
	status := attributes["status"].(string)

	i := personamodels.Inquiry{}

	createdAtStr := attributes["created-at"]
	if (createdAtStr != nil) {
		createdAt, _ := time.Parse(time.RFC3339, createdAtStr.(string))
		i.CreatedAt = &createdAt
	}
	if (attributes["started-at"] != nil) {
		startedAt, _ := time.Parse(time.RFC3339, attributes["started-at"].(string))
		i.StartedAt = &startedAt
	}
	if (attributes["completed-at"] != nil) {
		completedAt, _ := time.Parse(time.RFC3339, attributes["completed-at"].(string))
		i.CompletedAt = &completedAt
	}
	if (attributes["decisioned-at"] != nil) {
		decisionedAt, _ := time.Parse(time.RFC3339, attributes["decisioned-at"].(string))
		i.DecisionedAt = &decisionedAt
	}

	relationships := data["relationships"].(map[string]interface{})
	account := relationships["account"].(map[string]interface{})["data"].(map[string]interface{})
	acctId := account["id"].(string)

	i.PersInquiryID = inqId
	i.InquiryStatus = status
	i.PersAccountID = acctId
	
	return &i, nil
}

func UpdateInquiry(sqlDB *sql.DB, inquiry personamodels.Inquiry) {
	p := pers_db.PersDB{DB: sqlDB}
	p.InsertInquiry(inquiry)
}