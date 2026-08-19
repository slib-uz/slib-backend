package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func DraftEntityToModel(draft *entity.DraftEntity) *models.DraftModel {
	return models.NewDraftModel(
		draft.UserID,
		draft.Key,
		ToGormJson(draft.Data),
	)
}

func DraftModelToEntity(draftModel *models.DraftModel) *entity.DraftEntity {
	return entity.NewDraftEntity(
		draftModel.ID,
		draftModel.UserID,
		draftModel.Key,
		FromGormJson[map[string]any](draftModel.Data),
	)
}
