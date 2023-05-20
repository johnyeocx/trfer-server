package sumsubmodels

type IdDoc struct {
	IdDocType    string `json:"idDocType,omitempty"`
	Country      string `json:"country,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	FirstNameEn  string `json:"firstNameEn,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	MiddleNameEn string `json:"middleNameEn,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	LastNameEn   string `json:"lastNameEn,omitempty"`
	DateOfBirth  string `json:"dob,omitempty"` // yyyy-mm-dd format
}

type Info struct {
	FirstName    string  `json:"firstName,omitempty"`
	FirstNameEn  string  `json:"firstNameEn,omitempty"`
	MiddleName   string  `json:"middleName,omitempty"`
	MiddleNameEn string  `json:"middleNameEn,omitempty"`
	LastName     string  `json:"lastName,omitempty"`
	LastNameEn   string  `json:"lastNameEn,omitempty"`
	Dob          string  `json:"dob,omitempty"` //yyyy-mm-dd format
	Gender       string  `json:"gender,omitempty"`
	Country      string  `json:"country,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	IdDocs       []IdDoc `json:"idDocs,omitempty"`
}

type AccessToken struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

type Applicant struct {
	ID             string `json:"id,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	Key            string `json:"key,omitempty"`
	ClientID       string `json:"clientId,omitempty"`
	InspectionID   string `json:"inspectionId,omitempty"`
	ExternalUserID string `json:"externalUserId,omitempty"`
	Info           Info   `json:"info,omitempty"`
	FixedInfo      Info   `json:"fixedInfo,omitempty"`
	Review         struct {
		ElapsedSincePendingMs int    `json:"elapsedSincePendingMs,omitempty"`
		ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs,omitempty"`
		Reprocessing          bool   `json:"reprocessing,omitempty"`
		CreateDate            string `json:"createDate,omitempty"`
		ReviewDate            string `json:"reviewDate,omitempty"`
		StartDate             string `json:"startDate,omitempty"`
		ReviewResult          struct {
			ReviewAnswer string `json:"reviewAnswer,omitempty"`
		} `json:"reviewResult,omitempty"`
		ReviewStatus           string `json:"reviewStatus,omitempty"`
		NotificationFailureCnt int    `json:"notificationFailureCnt,omitempty"`
		Priority               int    `json:"priority,omitempty"`
	} `json:"review,omitempty"`
	Lang string `json:"lang,omitempty"`
	Type string `json:"type,omitempty"`
}