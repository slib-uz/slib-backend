package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AcademicDegreeModelToEntity(model *models.AcademicDegreeModel) *entity.AcademicDegreeEntity {
	return entity.NewAcademicDegreeEntity(
		model.ID,
		model.SourceID,
		model.Speciality,
		model.ConfirmedDate,
		model.UserID,
		UserModelToEntity(model.User),
		model.DiplomaNumber,
		model.ScienceSector,
		model.DegreeName,
		model.DegreeCode,
		model.ScienceSectorCode,
		model.Theme,
		model.AwardedAt,
	)
}

func AcademicDegreeEntityToModel(entity *entity.AcademicDegreeEntity) *models.AcademicDegreeModel {
	return models.NewAcademicDegreeModel(
		entity.SourceID,
		entity.Speciality,
		entity.ConfirmedDate,
		entity.UserID,
		entity.DiplomaNumber,
		entity.ScienceSector,
		entity.DegreeName,
		entity.DegreeCode,
		entity.ScienceSectorCode,
		entity.Theme,
		entity.AwardedAt,
	)
}
