package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalApplicationModelToEntity(model *models.JournalApplicationModel) *entity2.JournalApplicationEntity {

	var journal *entity2.JournalEntity
	var user *entity2.UserSharedEntity

	if model.Journal != nil {
		journal = JournalModelToEntity(model.Journal)
	}

	if model.User != nil {
		user = UserModelToSharedEntity(model.User)
	}

	return entity2.NewJournalApplicationEntity(
		model.ID,
		model.JournalID,
		journal,
		model.UserID,
		user,
		model.Status,
		model.RejectionReason,
		model.ReviewedAt,
	)

}
