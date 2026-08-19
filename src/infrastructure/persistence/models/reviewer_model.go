package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReviewerModel struct {
	gorm.Model
	ExternalID  uint   `gorm:"uniqueIndex;not null"`
	ScienceID   string `gorm:"uniqueIndex;not null"`
	FullName    string `gorm:"not null;size:256"`
	PhoneNumber string `gorm:"size:32"`
	Subjects    datatypes.JSON
}

func NewReviewerModel(externalID uint, scienceID string, fullName string, phoneNumber string, subjects datatypes.JSON) *ReviewerModel {
	return &ReviewerModel{ExternalID: externalID, ScienceID: scienceID, FullName: fullName, PhoneNumber: phoneNumber, Subjects: subjects}
}

func (ReviewerModel) TableName() string {
	return "reviewers"
}
