package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func AuthorshipClaimModelToEntity(model *models.AuthorshipClaimModel) *entity.AuthorshipClaimEntity {
	if model == nil {
		return nil
	}

	entity := entity.NewAuthorshipClaimEntity(
		model.ID,
		model.SenderID,
		UserModelToEntity(&model.Sender),
		model.ArticleID,
		ArticleModelToEntity(&model.Article),
		model.Comment,
		model.Status,
		model.CreatedAt,
		model.UpdatedAt,
	)

	entity.RejectReason = model.RejectReason
	entity.ReviewedByID = model.ReviewedByID
	entity.ReviewedAt = model.ReviewedAt

	if model.ReviewedBy != nil {
		entity.ReviewedBy = UserModelToEntity(model.ReviewedBy)
	}

	return entity
}

func AuthorshipClaimModelListToEntityList(models []*models.AuthorshipClaimModel) []*entity.AuthorshipClaimEntity {
	entities := make([]*entity.AuthorshipClaimEntity, len(models))
	for i, model := range models {
		entities[i] = AuthorshipClaimModelToEntity(model)
	}
	return entities
}

func AuthorshipClaimEntityToModel(entity *entity.AuthorshipClaimEntity) *models.AuthorshipClaimModel {
	if entity == nil {
		return nil
	}
	model := &models.AuthorshipClaimModel{
		SenderID:     entity.SenderID,
		ArticleID:    entity.ArticleID,
		Comment:      entity.Comment,
		Status:       entity.Status,
		RejectReason: entity.RejectReason,
		ReviewedByID: entity.ReviewedByID,
		ReviewedAt:   entity.ReviewedAt,
	}
	if entity.ID > 0 {
		model.ID = entity.ID
	}
	return model
}
