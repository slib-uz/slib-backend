package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type RoleModel struct {
	gorm.Model

	UserID uint       `gorm:"not null;index"`
	User   *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	PublisherID *uint           `gorm:"index"`
	Publisher   *PublisherModel `gorm:"foreignKey:PublisherID;constraint:OnDelete:CASCADE"`

	JournalID *uint         `gorm:"index"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE"`

	InstitutionID *uint              `gorm:"index"`
	Institution   *InstitutionModel  `gorm:"foreignKey:InstitutionID;constraint:OnDelete:CASCADE"`

	Role enum.UserRole `gorm:"default:0"`
	URL  *string       `gorm:"size:128"`
}

func (RoleModel) TableName() string {
	return "roles"
}

func NewUserRoleModel(userID uint, publisherID, journalID, institutionID *uint, role enum.UserRole, url *string) *RoleModel {
	return &RoleModel{
		UserID:        userID,
		PublisherID:   publisherID,
		JournalID:     journalID,
		InstitutionID: institutionID,
		Role:          role,
		URL:           url,
	}
}
