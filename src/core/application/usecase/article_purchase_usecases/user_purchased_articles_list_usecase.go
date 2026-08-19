package article_purchase_usecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserPurchasedArticlesListUsecase struct {
	repository repository.ArticlePurchaseRepository
}

// @inject
func NewUserPurchasedArticlesListUsecase(repository repository.ArticlePurchaseRepository) *UserPurchasedArticlesListUsecase {
	return &UserPurchasedArticlesListUsecase{repository: repository}
}

func (this UserPurchasedArticlesListUsecase) Execute(userID uint, page, pageSize int, search string) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error) {
	paging, err := this.repository.GetByUserID(userID, page, pageSize, search)
	if err != nil {
		return nil, err
	}
	return paging, nil
}
