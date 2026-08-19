package permissionusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type ArticlePurchasedPermission struct {
	repository repository.ArticlePurchaseRepository
}

// @inject
func NewArticlePurchasedPermission(repository repository.ArticlePurchaseRepository) *ArticlePurchasedPermission {
	return &ArticlePurchasedPermission{repository: repository}
}

func (this ArticlePurchasedPermission) Execute(articleID, userID uint) (bool, error) {
	return this.repository.IsExists(articleID, userID)
}
