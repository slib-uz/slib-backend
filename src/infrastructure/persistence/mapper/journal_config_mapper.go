package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalConfigEntityToModel(it *entity.JournalConfigEntity) *models.JournalConfigModel {
	return models.NewJournalConfigModel(
		it.JournalID,
		it.CreatorID,
		it.WebsiteURL,
		ToGormJson(it.Conf),
		it.IsActive,
	)
}

func JournalConfigModelToEntity(it *models.JournalConfigModel) *entity.JournalConfigEntity {
	var journal *entity.JournalEntity
	if it.Journal != nil {
		journal = JournalModelToEntity(it.Journal)
	}

	var creator *entity.UserEntity
	if it.Creator != nil {
		creator = UserModelToEntity(it.Creator)
		return entity.NewJournalConfigEntity(
			it.ID,
			it.JournalID,
			it.CreatorID,
			creator,
			journal,
			it.WebsiteURL,
			FromGormJson[map[string]any](it.Conf),
			it.IsActive,
			it.CreatedAt,
			it.UpdatedAt,
		)
	}

	return entity.NewJournalConfigEntity(
		it.ID,
		it.JournalID,
		it.CreatorID,
		creator,
		journal,
		it.WebsiteURL,
		FromGormJson[map[string]any](it.Conf),
		it.IsActive,
		it.CreatedAt,
		it.UpdatedAt,
	)
}
