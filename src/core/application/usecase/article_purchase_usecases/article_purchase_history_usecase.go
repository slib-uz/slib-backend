package article_purchase_usecases

import (
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticlePurchaseHistoryUseCase struct {
	repository               repository.ArticlePurchaseRepository
	publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase
}

// @inject
func NewArticlePurchaseHistoryUseCase(repository repository.ArticlePurchaseRepository, publisherAdminPermission *permissionusecases.PublisherAdminPermissionUseCase) *ArticlePurchaseHistoryUseCase {
	return &ArticlePurchaseHistoryUseCase{repository: repository, publisherAdminPermission: publisherAdminPermission}
}

func (this *ArticlePurchaseHistoryUseCase) Execute(user *entity2.UserBasicEntity, journalID uint, page, pageSize int, startDate, endDate time.Time) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error) {

	if !this.checkAccess(user, journalID) {
		return nil, response.PermissionDeniedError
	}

	paging, err := this.repository.GetByJournalID(journalID, page, pageSize, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return paging, nil
}

func (this *ArticlePurchaseHistoryUseCase) checkAccess(user *entity2.UserBasicEntity, journalID uint) bool {
	return user.IsAdmin || this.publisherAdminPermission.Execute(user.Roles, journalID)
}
