package article_purchase_usecases

import (
	"time"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticlePurchaseStatsUseCase struct {
	repository               repository.ArticlePurchaseRepository
	publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase
	journalRepository        repository.JournalRepository
}

// @inject
func NewArticlePurchaseStatsUseCase(repository repository.ArticlePurchaseRepository, publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase, journalRepository repository.JournalRepository) *ArticlePurchaseStatsUseCase {
	return &ArticlePurchaseStatsUseCase{repository: repository, publisherAdminPermission: publisherAdminPermission, journalRepository: journalRepository}
}

func (this *ArticlePurchaseStatsUseCase) Execute(user *entity2.UserBasicEntity, publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalArticlePurchaseStatsWithNameEntity], error) {

	if !this.checkAccess(user, publisherID) {
		return nil, response.PermissionDeniedError
	}

	paging, err := this.repository.StatsByJournal(publisherID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(paging.Page, paging.Size, paging.Total, mapper.JournalArticlePurchaseStatsEntityMapper(paging.Items)), nil
}

func (this *ArticlePurchaseStatsUseCase) checkAccess(user *entity2.UserBasicEntity, publisherID uint) bool {
	return user.IsAdmin || this.publisherAdminPermission.Execute(user.Roles, publisherID)
}
