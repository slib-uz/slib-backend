package journalconfigusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type CreateOrUpdateJournalConfigUseCase struct {
	repository        repository.JournalConfigRepository
	journalRepository repository.JournalRepository
}

// @inject
func NewCreateOrUpdateJournalConfigUseCase(repository repository.JournalConfigRepository, journalRepository repository.JournalRepository) *CreateOrUpdateJournalConfigUseCase {
	return &CreateOrUpdateJournalConfigUseCase{repository: repository, journalRepository: journalRepository}
}

func (this *CreateOrUpdateJournalConfigUseCase) Execute(c *entity.JournalConfigEntity, user *entity.UserBasicEntity) error {
	publisherID, err := this.journalRepository.GetPublisherIdByJournalId(c.JournalID)
	if err != nil {
		return err
	}

	if !this.hasUserPermission(user, c, publisherID) {
		return response.NewFailResponse(403, "You do not have permission to create or update JournalConfig for this journal.")
	}
	return this.repository.CreateOrUpdate(c)
}

func (this *CreateOrUpdateJournalConfigUseCase) hasUserPermission(user *entity.UserBasicEntity, c *entity.JournalConfigEntity, publisherID uint) bool {
	for _, r := range user.Roles {
		if r == nil {
			continue
		}

		if r.Role == enum.RoleAdmin {
			return true
		}

		if r.Role == enum.RolePublisherAdmin && r.PublisherID != nil && *r.PublisherID == publisherID {
			return true
		}

		if (r.Role == enum.RoleChiefEditor || r.Role == enum.RoleSecretary) &&
			r.JournalID != nil && *r.JournalID == c.JournalID {
			return true
		}
	}

	return false
}
