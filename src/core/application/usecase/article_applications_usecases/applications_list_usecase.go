package usecase

import (
	"time"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	permissionusecases2 "slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationsListUseCase struct {
	repository               repository.ApplicationRepository
	journalMemberPermission  *permissionusecases2.JournalMemberPermissionUseCase
	publisherAdminPermission *permissionusecases2.PublisherAdminPermissionUseCase
	journalRepository        repository.JournalRepository
}

// @inject
func NewApplicationsListUseCase(repository repository.ApplicationRepository, journalMemberPermission *permissionusecases2.JournalMemberPermissionUseCase, publisherAdminPermission *permissionusecases2.PublisherAdminPermissionUseCase, journalRepository repository.JournalRepository) *ApplicationsListUseCase {
	return &ApplicationsListUseCase{repository: repository, journalMemberPermission: journalMemberPermission, publisherAdminPermission: publisherAdminPermission, journalRepository: journalRepository}
}

func (this *ApplicationsListUseCase) Execute(user *entity2.UserBasicEntity, journalId uint, page, pageSize int, startDate, endDate time.Time, search string, status int) (*entity2.PagingEntity[entity2.ApplicationBasicEntity], error) {

	// If journalId is 0 (not provided), only admins can access all applications
	if journalId == 0 {
		if !user.IsAdmin {
			return nil, response.PermissionDeniedError
		}
	} else {
		// If journalId is provided, check access to that specific journal
		if err := this.checkAccess(user, journalId); err != nil {
			return nil, err
		}
	}

	paging, err := this.repository.GetByJournalID(journalId, page, pageSize, startDate, endDate, search, status)
	if err != nil {
		return nil, err
	}

	applications := make([]*entity2.ApplicationBasicEntity, len(paging.Items))
	for i, application := range paging.Items {
		applications[i] = mapper.ApplicationEntityToBasicEntity(application)
	}
	return entity2.NewPagingEntity(
		paging.Page,
		paging.Size,
		paging.Total,
		applications,
	), nil
}

func (this *ApplicationsListUseCase) checkAccess(user *entity2.UserBasicEntity, journalID uint) error {

	if user.IsAdmin {
		return nil
	}

	isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
	if err != nil {
		return err
	}
	if isMember {
		return nil
	}

	journal, err := this.journalRepository.FindByID(journalID)
	if err != nil {
		return err
	}

	if this.publisherAdminPermission.Execute(user.Roles, journal.PublisherID) {
		return nil
	}
	return response.PermissionDeniedError
}
