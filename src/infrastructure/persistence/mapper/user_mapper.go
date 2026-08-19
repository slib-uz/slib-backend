package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

func UserModelToEntity(model *models.UserModel) *entity.UserEntity {
	if model == nil {
		return nil
	}
	var roles []*entity.UserRoleEntity
	for _, role := range model.Roles {
		roles = append(roles, mapper.UserRoleModelToEntity(role))
	}

	return entity.NewUserEntity(
		model.ID,
		model.ScienceID,
		model.Pin,
		model.FirstName,
		model.LastName,
		model.MiddleName,
		model.FullName,
		model.PhoneNumber,
		model.IsAdmin,
		model.Gender,
		model.BirthDate,
		roles,
		model.Photo,
		model.Email,
		model.City,
		model.AcademicDegree,
		model.AcademicTitle,
		model.ORCIDID,
		0,
	)

}

func UserEntityToModel(entity *entity.UserEntity) *models.UserModel {
	return models.NewUserModel(
		entity.ScienceID,
		entity.Pin,
		entity.FirstName,
		entity.LastName,
		entity.MiddleName,
		entity.FullName,
		entity.Gender,
		entity.BirthDate,
		entity.PhoneNumber,
		entity.Photo,
		entity.Email,
		entity.City,
		entity.AcademicDegree,
		entity.AcademicTitle,
		entity.ORCIDID,
	)
}

func UserModelToSharedEntity(model *models.UserModel) *entity.UserSharedEntity {
	return entity.NewUserSharedEntity(
		model.ID,
		model.ScienceID,
		model.FullName,
	)
}
func UserModelListToEntityList(models []*models.UserModel) []*entity.UserEntity {
	entities := make([]*entity.UserEntity, len(models))
	for i, model := range models {
		entities[i] = UserModelToEntity(model)
	}
	return entities
}
