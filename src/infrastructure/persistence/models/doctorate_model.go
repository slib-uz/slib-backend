package models

import "gorm.io/gorm"

// from daraja.ilmiy.uz

type DoctorateModel struct {
	gorm.Model

	ExternalID         uint    `gorm:"uniqueIndex;not null"`
	DcType             *string `gorm:"size:64"`
	EduLang            string  `gorm:"size:32;not null"`
	Status             string  `gorm:"size:64;"`
	StatusCode         int
	AdmissionYear      *uint
	DirectionName      string  `gorm:"size:255;not null"`
	DirectionCode      string  `gorm:"size:16;not null"`
	AdvisorFullName    *string `gorm:"size:255"`
	AdvisorPin         *string `gorm:"size:32"`
	ScientificWorkName *string `gorm:"size:2048"`
	OrganizationTin    *string `gorm:"size:255;not null"`

	OrganizationID *uint              `gorm:"index"`
	Organization   *OrganizationModel `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"` // many2one
	UserID         *uint              `gorm:"not null;index"`
	User           *UserModel         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // many2one
}

func NewDoctorateModel(
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
	//organizationID *uint,
	userID *uint,
) *DoctorateModel {
	return &DoctorateModel{
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
		//OrganizationID:     organizationID,
		UserID: userID,
	}
}

func (DoctorateModel) TableName() string {
	return "doctorates"
}
