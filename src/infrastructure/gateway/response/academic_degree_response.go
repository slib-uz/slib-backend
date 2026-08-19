package response

import (
	"time"

	"slib.uz/src/core/utils"
)

type AcademicDegreeResponse struct {
	SourceID          uint                `json:"id"`
	Speciality        string              `json:"speciality"`
	ConfirmedDate     *utils.DateOnlyType `json:"confirmed_date"`
	DiplomaNumber     string              `json:"diploma_number"`
	ScienceSector     string              `json:"science_sector"`
	DegreeName        string              `json:"degree_name"`
	DegreeCode        int                 `json:"degree_code"`
	ScienceSectorCode int                 `json:"science_sector_code"`
	Theme             string              `json:"theme"`
}

func NewAcademicDegreeResponse(sourceID uint, speciality string, confirmedDate *time.Time, diplomaNumber string, scienceSector string, degreeName string, degreeCode int, scienceSectorCode int, theme string) *AcademicDegreeResponse {
	var dateOnly *utils.DateOnlyType
	if confirmedDate != nil {
		dateOnly = &utils.DateOnlyType{Time: *confirmedDate}
	}
	return &AcademicDegreeResponse{
		SourceID:          sourceID,
		Speciality:        speciality,
		ConfirmedDate:     dateOnly,
		DiplomaNumber:     diplomaNumber,
		ScienceSector:     scienceSector,
		DegreeName:        degreeName,
		DegreeCode:        degreeCode,
		ScienceSectorCode: scienceSectorCode,
		Theme:             theme,
	}
}
