package articleusecases

import (
	"fmt"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/application/validation"
	"slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type ArticleUpdateUseCase struct {
	articleRepository             repository.ArticleRepository
	publishedArticleRepository    repository.PublishedArticleRepository
	languageRepository            repository.LanguageRepository
	studyFieldRepository          repository.StudyFieldRepository
	authorRepository              repository.AuthorRepository
	storage                       storage.FileStorage
	hasAnyJournalRolePermission   *permissionusecases.HasAnyJournalRolePermissionUseCase
	hasAnyPublisherRolePermission *permissionusecases.HasAnyPublisherRolePermissionUseCase
	updateRoiSenderTask           *tasks.UpdateRoiSenderTask
}

// @inject
func NewArticleUpdateUseCase(
	articleRepository repository.ArticleRepository,
	publishedArticleRepository repository.PublishedArticleRepository,
	languageRepository repository.LanguageRepository,
	studyFieldRepository repository.StudyFieldRepository,
	authorRepository repository.AuthorRepository,
	storage storage.FileStorage,
	hasAnyJournalRolePermission *permissionusecases.HasAnyJournalRolePermissionUseCase,
	hasAnyPublisherRolePermission *permissionusecases.HasAnyPublisherRolePermissionUseCase,
	updateRoiSenderTask *tasks.UpdateRoiSenderTask,
) *ArticleUpdateUseCase {
	return &ArticleUpdateUseCase{
		articleRepository:             articleRepository,
		publishedArticleRepository:    publishedArticleRepository,
		languageRepository:            languageRepository,
		studyFieldRepository:          studyFieldRepository,
		authorRepository:              authorRepository,
		storage:                       storage,
		hasAnyJournalRolePermission:   hasAnyJournalRolePermission,
		hasAnyPublisherRolePermission: hasAnyPublisherRolePermission,
		updateRoiSenderTask:           updateRoiSenderTask,
	}
}

func (this *ArticleUpdateUseCase) Execute(user *entity.UserBasicEntity, articleID uint, e *entity.ArticleManageUpdateEntity) error {
	article, err := this.articleRepository.GetByIDWithJournal(articleID)
	if err != nil {
		return err
	}

	if err := this.checkAccess(user, article); err != nil {
		return err
	}

	if e.LanguageID != 0 {
		if err := this.checkId("Language", e.LanguageID, this.languageRepository.IsExist); err != nil {
			return err
		}
	}

	if err := validation.ValidateIDs("StudyField", e.StudyFieldsIDs, this.studyFieldRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("Author", e.CoAuthorsIDs, this.authorRepository.ExistingIds); err != nil {
		return err
	}

	//if e.ContentFile != "" && e.ContentFile != article.ContentFile {
	//	contentPath, err := this.storage.UploadIfExists(enum.FolderArticles, e.ContentFile, enum.BucketPrivate)
	//	if err != nil {
	//		return err
	//	}
	//	e.ContentFile = contentPath
	//}
	//
	//if e.ExpertConclusionFile != nil && *e.ExpertConclusionFile != "" && *e.ExpertConclusionFile != *article.ExpertConclusionFile {
	//	expertPath, err := this.storage.UploadIfExists(enum.FolderArticles, *e.ExpertConclusionFile, enum.BucketPrivate)
	//	if err != nil {
	//		return err
	//	}
	//	e.ExpertConclusionFile = &expertPath
	//}

	err = this.articleRepository.ManageUpdate(articleID, e)
	if err != nil {
		return err
	}

	publishedArticle, err := this.publishedArticleRepository.GetByIDWithRelations(articleID)
	if err != nil {
		return err
	}

	// UpdateRoiSenderTask.Run navbat xatosini o'zi loglaydi va har doim
	// nil qaytaradi — maqolani yangilash ROI navbatiga bog'liq emas.
	_ = this.sendToROI(publishedArticle)

	return nil
}

func (this *ArticleUpdateUseCase) checkAccess(user *entity.UserBasicEntity, article *entity.ArticleBasicEntity) error {
	if user.IsAdmin {
		return nil
	}

	if article.Journal != nil {
		if this.hasAnyPublisherRolePermission.Execute(user.Roles, article.Journal.PublisherID) {
			return nil
		}
	}

	if this.hasAnyJournalRolePermission.Execute(user.Roles, article.JournalID) {
		return nil
	}

	return response.PermissionDeniedError
}

func (this *ArticleUpdateUseCase) checkId(field string, id uint, exists func(id uint) (bool, error)) error {
	ex, err := exists(id)
	if err != nil {
		return err
	}
	if !ex {
		return response.NewFailResponse(404, fmt.Sprintf("%s not found", field))
	}
	return nil
}

func (this *ArticleUpdateUseCase) sendToROI(article *entity.ArticleEntity) error {
	return this.updateRoiSenderTask.Run(tasks.UpdateRoiSenderPayload{
		ArticleID: article.ID,
	})
}
