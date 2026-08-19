package models

import (
	"gorm.io/gorm"
	"time"
)

type AcademicTitleModel struct {
	gorm.Model

	SourceID uint `gorm:"not null;index"`

	Title             string
	ConfirmedDate     *time.Time `gorm:"type:date"`
	DiplomaNumber     string
	UserID            uint       `gorm:"not null"`
	User              *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ScienceSector     string
	ScienceSectorCode int
	TitleCode         int
	Speciality        string
	AwardedAt         *time.Time `gorm:"type:date;not null"`
}

func (AcademicTitleModel) TableName() string {
	return "academic_titles"
}

func NewAcademicTitleModel(sourceID uint, title string, confirmedDate *time.Time, userID uint, diplomaNumber string, scienceSector string, scienceSectorCode int, titleCode int, speciality string, awardedAt *time.Time) *AcademicTitleModel {
	return &AcademicTitleModel{
		SourceID:          sourceID,
		Title:             title,
		ConfirmedDate:     confirmedDate,
		UserID:            userID,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		ScienceSectorCode: scienceSectorCode,
		TitleCode:         titleCode,
		Speciality:        speciality,
		AwardedAt:         awardedAt,
	}
}
