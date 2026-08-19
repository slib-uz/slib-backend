package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AuthorModelToEntity(model *models.AuthorModel) *entity2.AuthorEntity {
	var user *entity2.UserEntity
	var orcidID *string
	if model.Users != nil && len(*model.Users) > 0 {
		first := &(*model.Users)[0]
		user = UserModelToEntity(first)
		orcidID = first.ORCIDID
	}

	return entity2.NewAuthorEntity(
		model.ID,
		model.FullName,
		model.ScienceID,
		0,
		nil,
		nil,
		nil,
		nil,
		orcidID,
		user,
	)
}

func AuthorEntityToModel(entity *entity2.AuthorEntity) models.AuthorModel {

	return models.AuthorModel{
		FullName:  entity.FullName,
		ScienceID: entity.ScienceID,
	}
}
