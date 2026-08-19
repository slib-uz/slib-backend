package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func AcademicDegreeResToEntity(res *response.AcademicDegreeResponse) *entity.AcademicDegreeEntity {
	if res == nil {
		return nil
	}

	var confirmedDate *time.Time
	if res.ConfirmedDate != nil {
		t := res.ConfirmedDate.Time
		confirmedDate = &t
	}

	return entity.NewAcademicDegreeEntity(
		0,
		res.SourceID,
		res.Speciality,
		confirmedDate,
		0,
		nil,
		res.DiplomaNumber,
		res.ScienceSector,
		res.DegreeName,
		res.DegreeCode,
		res.ScienceSectorCode,
		res.Theme,
		nil,
	)
}
