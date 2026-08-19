package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ReviewerEntityToModel(reviewer *entity.ReviewerEntity) *models.ReviewerModel {
	return models.NewReviewerModel(
		reviewer.ExternalID,
		reviewer.ScienceID,
		reviewer.FullName,
		reviewer.PhoneNumber,
		ToGormJson(reviewer.Subjects),
	)
}

func ReviewerModelToEntity(reviewer *models.ReviewerModel) *entity.ReviewerEntity {
	return entity.NewReviewerEntity(
		reviewer.ID,
		reviewer.ExternalID,
		reviewer.ScienceID,
		reviewer.FullName,
		reviewer.PhoneNumber,
		FromGormJson[[]map[string]any](reviewer.Subjects),
	)
}
