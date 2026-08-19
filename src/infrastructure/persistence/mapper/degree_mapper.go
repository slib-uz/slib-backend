package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func DegreeModelToEntity(degree *models.DegreeModel) *entity.DegreeEntity {
	return entity.NewDegreeEntity(
		degree.ID,
		degree.DegreeTypeID,
		degree.DegreeType,
		degree.Field,
		degree.DegreeStatusID,
		degree.DegreeStatusName,
		degree.ConfirmedDate,
		degree.Protocol,
		degree.UserID,
	)
}

func DegreeEntityToModel(userID uint, degree *entity.DegreeEntity) *models.DegreeModel {
	return &models.DegreeModel{
		DegreeTypeID:     degree.DegreeTypeID,
		DegreeType:       degree.DegreeType,
		Field:            degree.Field,
		DegreeStatusID:   degree.DegreeStatusID,
		DegreeStatusName: degree.DegreeStatusName,
		ConfirmedDate:    degree.ConfirmedDate,
		Protocol:         degree.Protocol,
		UserID:           &userID,
	}
}

func DegreeUpdateMapper(existing, degree *models.DegreeModel) *models.DegreeModel {
	existing.DegreeTypeID = degree.DegreeTypeID
	existing.DegreeType = degree.DegreeType
	existing.Field = degree.Field
	existing.DegreeStatusID = degree.DegreeStatusID
	existing.DegreeStatusName = degree.DegreeStatusName
	existing.ConfirmedDate = degree.ConfirmedDate
	existing.Protocol = degree.Protocol
	return existing
}
