package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ArticleAuthorAffiliationListModelToEntity(models []*models.ArticleAuthorAffiliationModel) []*entity.ArticleAuthorAffiliationEntity {
	entities := make([]*entity.ArticleAuthorAffiliationEntity, len(models))
	for i, model := range models {
		entities[i] = ArticleAuthorAffiliationModelToEntity(model)
	}
	return entities
}

func ArticleAuthorAffiliationListEntityToModel(entities []*entity.ArticleAuthorAffiliationEntity) []*models.ArticleAuthorAffiliationModel {
	models := make([]*models.ArticleAuthorAffiliationModel, len(entities))
	for i, entity := range entities {
		models[i] = ArticleAuthorAffiliationEntityToModel(entity)
	}
	return models
}

func ArticleAuthorAffiliationEntityToModel(entity *entity.ArticleAuthorAffiliationEntity) *models.ArticleAuthorAffiliationModel {
	return &models.ArticleAuthorAffiliationModel{
		ArticleID:        entity.ArticleID,
		AuthorID:         entity.AuthorID,
		OrganizationID:   entity.OrganizationID,
		OrganizationName: entity.OrganizationName,
		OrganizationTin:  entity.OrganizationTin,
		PositionName:     entity.PositionName,
	}
}

func ArticleAuthorAffiliationModelToEntity(model *models.ArticleAuthorAffiliationModel) *entity.ArticleAuthorAffiliationEntity {
	var author *entity.AuthorEntity
	var scienceID string
	if model.Author != nil {
		author = AuthorModelToEntity(model.Author)
		scienceID = author.ScienceID
	}
	return &entity.ArticleAuthorAffiliationEntity{
		ID:               model.ID,
		ArticleID:        model.ArticleID,
		AuthorID:         model.AuthorID,
		ScienceID:        scienceID,
		Author:           author,
		OrganizationID:   model.OrganizationID,
		OrganizationName: model.OrganizationName,
		OrganizationTin:  model.OrganizationTin,
		PositionName:     model.PositionName,
	}
}
