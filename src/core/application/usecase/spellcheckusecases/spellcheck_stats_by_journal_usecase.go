package spellcheckusecases

import (
	"time"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SpellCheckResultStatsByJournalUseCase struct {
	repository               repository.SpellCheckResultRepository
	publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase
}

// @inject
func NewSpellCheckResultStatsByJournalUseCase(repository repository.SpellCheckResultRepository, publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase) *SpellCheckResultStatsByJournalUseCase {
	return &SpellCheckResultStatsByJournalUseCase{repository: repository, publisherAdminPermission: publisherAdminPermission}
}

func (this *SpellCheckResultStatsByJournalUseCase) Execute(user *entity2.UserBasicEntity, publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalSpellcheckStatsWithNameEntity], error) {

	if err := this.checkAccess(user, publisherID); err != nil {
		return nil, err
	}

	paging, err := this.repository.StatsByJournal(publisherID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(paging.Page, paging.Size, paging.Total, mapper.JournalSpellcheckStatsMapper(paging.Items)), nil
}

func (this *SpellCheckResultStatsByJournalUseCase) checkAccess(user *entity2.UserBasicEntity, publisherID uint) error {
	if publisherID == 0 {
		if user.IsAdmin {
			return nil
		}
		return response.PermissionDeniedError
	}

	if this.publisherAdminPermission.Execute(user.Roles, publisherID) {
		return nil
	}

	return response.PermissionDeniedError

}
