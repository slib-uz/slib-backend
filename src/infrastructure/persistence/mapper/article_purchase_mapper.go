package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ArticlePurchaseModelToEntity(purchase *models.ArticlePurchaseModel) *entity.ArticlePurchaseEntity {
	var article *entity.ArticleBasicEntity
	var user *entity.UserSharedEntity
	var journal *entity.JournalBasicEntity

	if purchase.Article != nil {
		article = ArticleModelToBasicEntity(purchase.Article)
	}
	if purchase.User != nil {
		user = UserModelToSharedEntity(purchase.User)
	}

	if purchase.Journal != nil {
		journal = JournalModelToBasicEntity(purchase.Journal)
	}

	return entity.NewArticlePurchaseEntity(
		purchase.ID,
		purchase.JournalID,
		journal,
		purchase.ArticleID,
		article,
		purchase.UserID,
		user,
		purchase.Amount,
	)
}

func ArticlePurchaseEntityToModel(purchase *entity.ArticlePurchaseEntity) *models.ArticlePurchaseModel {
	return models.NewArticlePurchaseModel(
		purchase.JournalID,
		purchase.ArticleID,
		purchase.UserID,
		purchase.Amount,
	)
}
