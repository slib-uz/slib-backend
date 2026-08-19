package antiplagusecases

import (
	"time"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AntiPlagStatsByJournalUseCase struct {
	repository               repository.AntiPlagRepository
	publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase
}

// @inject
func NewAntiPlagStatsByJournalUseCase(repository repository.AntiPlagRepository, publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase) *AntiPlagStatsByJournalUseCase {
	return &AntiPlagStatsByJournalUseCase{repository: repository, publisherAdminPermission: publisherAdminPermission}
}

func (this *AntiPlagStatsByJournalUseCase) Execute(user *entity2.UserBasicEntity, publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalAntiPlagStatsWithNameEntity], error) {

	if err := this.checkAccess(user, publisherID); err != nil {
		return nil, err
	}

	paging, err := this.repository.StatsByJournal(publisherID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, err
	}

	mappedItems := mapper.JournalAntiPlagStatsMapper(paging.Items)

	return entity2.NewPagingEntity(paging.Page, paging.Size, paging.Total, mappedItems), nil
}

func (this *AntiPlagStatsByJournalUseCase) checkAccess(user *entity2.UserBasicEntity, publisherID uint) error {
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
