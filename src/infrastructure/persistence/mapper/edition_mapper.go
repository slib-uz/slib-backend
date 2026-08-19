package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func EditionModelToEntity(e *models.EditionModel) *entity.EditionEntity {
	var journal *entity.JournalEntity
	if e.Journal != nil {
		journal = JournalModelToEntity(e.Journal)
	}

	return entity.NewEditionEntity(
		e.ID,
		e.JournalID,
		journal,
		e.File,
		e.CoverImage,
		e.Name,
		e.Number,
		e.Volume,
		e.PublishedAt,
	)
}

func EditionEntityToModel(e *entity.EditionEntity) *models.EditionModel {
	return models.NewEditionModel(
		e.JournalID,
		e.File,
		e.CoverImage,
		e.Name,
		e.Number,
		e.Volume,
		e.PublishedAt,
	)
}
