package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AcademicTitleModelToEntity(model *models.AcademicTitleModel) *entity.AcademicTitleEntity {
	return entity.NewAcademicTitleEntity(
		model.ID,
		model.SourceID,
		model.Title,
		model.ConfirmedDate,
		model.UserID,
		UserModelToEntity(model.User),
		model.DiplomaNumber,
		model.ScienceSector,
		model.ScienceSectorCode,
		model.TitleCode,
		model.Speciality,
		model.AwardedAt,
	)
}

func AcademicTitleEntityToModel(entity *entity.AcademicTitleEntity) *models.AcademicTitleModel {
	return models.NewAcademicTitleModel(
		entity.SourceID,
		entity.Title,
		entity.ConfirmedDate,
		entity.UserID,
		entity.DiplomaNumber,
		entity.ScienceSector,
		entity.ScienceSectorCode,
		entity.TitleCode,
		entity.Speciality,
		entity.AwardedAt,
	)
}
