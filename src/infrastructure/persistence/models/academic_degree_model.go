package models

import (
	"gorm.io/gorm"
	"time"
)

type AcademicDegreeModel struct {
	gorm.Model

	SourceID          uint `gorm:"not null;index"`
	Speciality        string
	ConfirmedDate     *time.Time `gorm:"type:date"`
	UserID            uint       `gorm:"not null"`
	User              *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	DiplomaNumber     string
	ScienceSector     string
	DegreeName        string
	DegreeCode        int
	ScienceSectorCode int
	Theme             string
	AwardedAt         *time.Time `gorm:"type:date;not null"`
}

func (AcademicDegreeModel) TableName() string {
	return "academic_degrees"
}

func NewAcademicDegreeModel(sourceID uint, speciality string, confirmedDate *time.Time, userID uint, diplomaNumber string, scienceSector string, degreeName string, degreeCode int, scienceSectorCode int, theme string, awardedAt *time.Time) *AcademicDegreeModel {
	return &AcademicDegreeModel{
		SourceID:          sourceID,
		Speciality:        speciality,
		ConfirmedDate:     confirmedDate,
		UserID:            userID,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		DegreeName:        degreeName,
		DegreeCode:        degreeCode,
		ScienceSectorCode: scienceSectorCode,
		Theme:             theme,
		AwardedAt:         awardedAt,
	}
}
