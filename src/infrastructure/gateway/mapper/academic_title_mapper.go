package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func AcademicTitleResToEntity(res *response.AcademicTitleResponse) *entity.AcademicTitleEntity {
	if res == nil {
		return nil
	}

	var confirmedDate *time.Time
	if res.ConfirmedDate != nil {
		t := res.ConfirmedDate.Time
		confirmedDate = &t
	}

	return entity.NewAcademicTitleEntity(
		0,
		res.SourceID,
		res.Title,
		confirmedDate,
		0,
		nil,
		res.DiplomaNumber,
		res.ScienceSector,
		res.ScienceSectorCode,
		res.TitleCode,
		res.Speciality,
		nil,
	)
}
