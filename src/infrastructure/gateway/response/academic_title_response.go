package response

import (
	"time"

	"slib.uz/src/core/utils"
)

type AcademicTitleResponse struct {
	SourceID          uint                `json:"id"`
	Title             string              `json:"title"`
	ConfirmedDate     *utils.DateOnlyType `json:"confirmed_date"`
	DiplomaNumber     string              `json:"diploma_number"`
	ScienceSector     string              `json:"science_sector"`
	ScienceSectorCode int                 `json:"science_sector_code"`
	TitleCode         int                 `json:"title_code"`
	Speciality        string              `json:"speciality"`
}

func NewAcademicTitleResponse(sourceID uint, title string, confirmedDate *time.Time, diplomaNumber string, scienceSector string, scienceSectorCode int, titleCode int, speciality string) *AcademicTitleResponse {
	var dateOnly *utils.DateOnlyType
	if confirmedDate != nil {
		dateOnly = &utils.DateOnlyType{Time: *confirmedDate}
	}
	return &AcademicTitleResponse{
		SourceID:          sourceID,
		Title:             title,
		ConfirmedDate:     dateOnly,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		ScienceSectorCode: scienceSectorCode,
		TitleCode:         titleCode,
		Speciality:        speciality,
	}
}
