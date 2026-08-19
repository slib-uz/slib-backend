package uploadusecases

import (
	"time"

	"slib.uz/src/core/application/response"
	permissionusecases2 "slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type PresignedURLUseCase struct {
	articleRepository          repository.ArticleRepository
	storage                    storage.FileStorage
	journalMemberPermission    *permissionusecases2.JournalMemberPermissionUseCase
	articlePurchasedPermission *permissionusecases2.ArticlePurchasedPermission
	authorPermission           *permissionusecases2.ArticleAuthorPermissionUseCase
}

// @inject
func NewPresignedURLUseCase(articleRepository repository.ArticleRepository, storage storage.FileStorage, journalMemberPermission *permissionusecases2.JournalMemberPermissionUseCase, articlePurchasedPermission *permissionusecases2.ArticlePurchasedPermission, authorPermission *permissionusecases2.ArticleAuthorPermissionUseCase) *PresignedURLUseCase {
	return &PresignedURLUseCase{articleRepository: articleRepository, storage: storage, journalMemberPermission: journalMemberPermission, articlePurchasedPermission: articlePurchasedPermission, authorPermission: authorPermission}
}

func (this *PresignedURLUseCase) Execute(user *entity.UserBasicEntity, articleID uint) (string, error) {

	article, err := this.articleRepository.FindByID(articleID)
	if err != nil {
		return "", err
	}

	url, err := this.storage.PresignedURL(enum2.BucketPrivate, article.ContentFile, time.Minute*10)
	if err != nil {
		return "", err
	}

	if article.AccessType == enum2.PublicAccessType {
		return url, nil
	}

	if access, err := this.checkAccess(user, article.JournalID, articleID); err != nil {
		return "", err
	} else if !access {
		return "", response.PermissionDeniedError
	}

	return url, nil

}

func (this *PresignedURLUseCase) checkAccess(user *entity.UserBasicEntity, journalID, articleID uint) (bool, error) {

	isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
	if err != nil {
		return false, err
	}
	if isMember {
		return true, nil
	}

	if isAuthor, err := this.authorPermission.Execute(user.ScienceID, articleID); err != nil {
		return false, err
	} else if isAuthor {
		return true, nil
	}

	if isPurchased, err := this.articlePurchasedPermission.Execute(articleID, user.ID); err != nil {
		return false, err
	} else if isPurchased {
		return true, nil
	}

	return false, nil

}
