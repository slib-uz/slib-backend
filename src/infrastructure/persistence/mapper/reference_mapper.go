package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ReferenceModelToEntity(r *models.ReferenceModel) *entity2.ReferenceEntity {
	return entity2.NewReferenceEntity(
		r.ID,
		r.Name,
		r.ArticleID,
		nil, // Don't recursively load Article to avoid circular reference
	)
}

func ReferenceModelListToEntityList(refs []*models.ReferenceModel) []*entity2.ReferenceEntity {
	result := make([]*entity2.ReferenceEntity, len(refs))
	for i, r := range refs {
		result[i] = ReferenceModelToEntity(r)
	}
	return result
}
