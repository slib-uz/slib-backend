package entity

type DoctorateEntity struct {
	ID                 uint    `json:"id"`
	ExternalID         uint    `json:"external_id"`
	DcType             *string `json:"dc_type"`
	EduLang            string  `json:"edu_lang"`
	Status             string  `json:"status"`
	StatusCode         int     `json:"status_code"`
	AdmissionYear      *uint   `json:"admission_year"`
	DirectionName      string  `json:"direction_name"`
	DirectionCode      string  `json:"direction_code"`
	AdvisorFullName    *string `json:"advisor_full_name"`
	AdvisorPin         *string `json:"advisor_pin"`
	ScientificWorkName *string `json:"scientific_work_name"`
	OrganizationTin    *string `json:"organization_tin"`
	UserID             *uint   `json:"user_id"`
}

func NewDoctorateEntity(
	ID uint,
	externalID uint,
	dcType *string,
	eduLang string,
	status string,
	statusCode int,
	admissionYear *uint,
	directionName string,
	directionCode string,
	advisorFullName *string,
	advisorPin *string,
	scientificWorkName *string,
	organizationTin *string,
	userID *uint,
) *DoctorateEntity {
	return &DoctorateEntity{
		ID:                 ID,
		ExternalID:         externalID,
		DcType:             dcType,
		EduLang:            eduLang,
		Status:             status,
		StatusCode:         statusCode,
		AdmissionYear:      admissionYear,
		DirectionName:      directionName,
		DirectionCode:      directionCode,
		AdvisorFullName:    advisorFullName,
		AdvisorPin:         advisorPin,
		ScientificWorkName: scientificWorkName,
		OrganizationTin:    organizationTin,
		UserID:             userID,
	}
}
