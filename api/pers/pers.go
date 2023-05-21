package persona

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/johnyeocx/usual/server2/db/models"
	personamodels "github.com/johnyeocx/usual/server2/db/models/persona_models"
	"github.com/johnyeocx/usual/server2/db/pers_db"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/errors/pers_errors"
	"github.com/johnyeocx/usual/server2/errors/user_errors"
	"github.com/johnyeocx/usual/server2/persona"
	"github.com/johnyeocx/usual/server2/utils/enums/persona/InquiryStatus"
)

func GetInquiryAccessToken(sqlDB *sql.DB, uId int) (map[string]interface{}, *models.RequestError){
	// Get latest inquiry
	p := pers_db.PersDB{DB: sqlDB}
	inq, err := p.GetUserLatestInquiry(uId)
	if err != nil {
		return nil, pers_errors.GetUserInquiryFailedErr(err)
	}


	var sessionToken *string
	if (inq.InquiryStatus == string(InquiryStatus.Created))  {
		fmt.Println("Here")
		stoken, err := persona.GetInquirySessionToken(inq.PersInquiryID)
		fmt.Println("Stoken:", stoken)
		if err == nil {
			sessionToken = &stoken
		}
	}

	return map[string]interface{}{
		"session_token": sessionToken,
		"inquiry": inq,
	}, nil
}

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

func DecodeInquirySessionWebhook(data map[string]interface{}) (*personamodels.Inquiry, *models.RequestError){
	sessionId := data["id"].(string)

	attributes := data["attributes"].(map[string]interface{})
	sessionStatus := attributes["status"].(string)

	relationships := data["relationships"].(map[string]interface{})
	account := relationships["inquiry"].(map[string]interface{})["data"].(map[string]interface{})
	inqId := account["id"].(string)

	i := personamodels.Inquiry{}
	i.PersInquiryID = inqId
	i.PersSessionID = sessionId
	i.SessionStatus = sessionStatus
	fmt.Println("Inquiry Session: ", i)
	
	return &i, nil
}

func UpdateInquiry(sqlDB *sql.DB, inquiry personamodels.Inquiry) (*models.RequestError) {
	p := pers_db.PersDB{DB: sqlDB}
	err := p.InsertInquiry(inquiry)
	if err != nil {
		return pers_errors.UpdateInquiryFailedErr(err)
	}

	if inquiry.InquiryStatus == string(InquiryStatus.Approved) {
		u := user_db.UserDB{DB: sqlDB}
		err := u.SetPersApproved(inquiry.PersAccountID, true)
		if err != nil {
			return user_errors.SetPersApprovedFailedErr(err)
		}
	}

	return nil
}

func UpdateInquirySession(sqlDB *sql.DB, inquiry personamodels.Inquiry) (*models.RequestError) {
	p := pers_db.PersDB{DB: sqlDB}
	err := p.UpdateInquirySession(inquiry)
	if err != nil {
		return pers_errors.UpdateInquiryFailedErr(err)
	}

	return nil
}