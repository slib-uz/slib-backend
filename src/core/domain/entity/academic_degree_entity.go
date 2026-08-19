package entity

import "time"

type AcademicDegreeEntity struct {
	ID                uint        `json:"id"`
	SourceID          uint        `json:"source_id"`
	Speciality        string      `json:"speciality"`
	ConfirmedDate     *time.Time  `json:"confirmed_date,omitempty"`
	UserID            uint        `json:"user_id"`
	User              *UserEntity `json:"user"`
	DiplomaNumber     string      `json:"diploma_number"`
	ScienceSector     string      `json:"science_sector"`
	DegreeName        string      `json:"degree_name"`
	DegreeCode        int         `json:"degree_code"`
	ScienceSectorCode int         `json:"science_sector_code"`
	Theme             string      `json:"theme"`
	AwardedAt         *time.Time  `json:"awarded_at,omitempty"`
}

func NewAcademicDegreeEntity(id uint, sourceID uint, speciality string, confirmedDate *time.Time, userID uint, user *UserEntity, diplomaNumber string, scienceSector string, degreeName string, degreeCode int, scienceSectorCode int, theme string, awardedAt *time.Time) *AcademicDegreeEntity {
	return &AcademicDegreeEntity{
		ID:                id,
		SourceID:          sourceID,
		Speciality:        speciality,
		ConfirmedDate:     confirmedDate,
		UserID:            userID,
		User:              user,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		DegreeName:        degreeName,
		DegreeCode:        degreeCode,
		ScienceSectorCode: scienceSectorCode,
		Theme:             theme,
		AwardedAt:         awardedAt,
	}
}
