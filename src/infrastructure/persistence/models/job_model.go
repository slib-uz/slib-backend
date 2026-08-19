package models

import (
	"gorm.io/gorm"
)

type JobModel struct {
	gorm.Model

	OrganizationTin  string             `gorm:"size:32;index;not null'"`
	OrganizationName string             `gorm:"size:2048;index;"`
	OrganizationID   *uint              `gorm:"column:organization_id;index"`
	Organization     *OrganizationModel `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"` // many2one

	UserID uint       `gorm:"not null;index"`
	User   *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // many2one

	PositionName string `gorm:"size:100"`
}

func (JobModel) TableName() string {
	return "jobs"
}

func NewJobModel(organizationID *uint, userID uint, organizationTin, organizationName, positionName string) *JobModel {
	return &JobModel{
		OrganizationID:   organizationID,
		UserID:           userID,
		OrganizationTin:  organizationTin,
		OrganizationName: organizationName,
		PositionName:     positionName,
	}

}
