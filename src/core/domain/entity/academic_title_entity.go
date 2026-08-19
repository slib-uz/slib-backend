package entity

import "time"

type AcademicTitleEntity struct {
	ID                uint        `json:"id"`
	SourceID          uint        `json:"source_id"`
	Title             string      `json:"title"`
	ConfirmedDate     *time.Time  `json:"confirmed_date,omitempty"`
	UserID            uint        `json:"user_id"`
	User              *UserEntity `json:"user"`
	DiplomaNumber     string      `json:"diploma_number"`
	ScienceSector     string      `json:"science_sector"`
	ScienceSectorCode int         `json:"science_sector_code"`
	TitleCode         int         `json:"title_code"`
	Speciality        string      `json:"speciality"`
	AwardedAt         *time.Time  `json:"awarded_at,omitempty"`
}

func NewAcademicTitleEntity(id uint, sourceID uint, title string, confirmedDate *time.Time, userID uint, user *UserEntity, diplomaNumber string, scienceSector string, scienceSectorCode int, titleCode int, speciality string, awardedAt *time.Time) *AcademicTitleEntity {
	return &AcademicTitleEntity{
		ID:                id,
		SourceID:          sourceID,
		Title:             title,
		ConfirmedDate:     confirmedDate,
		UserID:            userID,
		User:              user,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		ScienceSectorCode: scienceSectorCode,
		TitleCode:         titleCode,
		Speciality:        speciality,
		AwardedAt:         awardedAt,
	}
}
