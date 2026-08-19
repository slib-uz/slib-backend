package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func SpellCheckResultEntityToResponseEntity(data *entity.SpellCheckResultEntity) *entity.SpellCheckResultEntity {
	return entity.NewSpellCheckResultEntity(
		data.ID,
		data.ReviewStageID,
		data.ReviewStage,
		data.ApplicationID,
		data.Application,
		data.JournalID,
		data.Journal,
		data.File,
		data.ResultFile,
		data.Status,
		data.SubmitterID,
		data.Submitter,
		data.ResultTime,
	)
}
