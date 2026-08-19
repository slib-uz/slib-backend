package articleusecases

import (
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/application/validation"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DirectArticleCreateUseCase struct {
	articleRepository                  repository.ArticleRepository
	studyFieldRepository               repository.StudyFieldRepository
	authorRepository                   repository.AuthorRepository
	languageRepository                 repository.LanguageRepository
	tagRepository                      repository.TagRepository
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository
	permissionUseCase                  *permissionusecases.JournalMemberPermissionUseCase
	publishRoiSenderTask               *tasks.PublishRoiSenderTask
}

// @inject
func NewDirectArticleCreateUseCase(
	articleRepository repository.ArticleRepository,
	studyFieldRepository repository.StudyFieldRepository,
	authorRepository repository.AuthorRepository,
	languageRepository repository.LanguageRepository,
	tagRepository repository.TagRepository,
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository,
	permissionUseCase *permissionusecases.JournalMemberPermissionUseCase,
	publishRoiSenderTask *tasks.PublishRoiSenderTask,
) *DirectArticleCreateUseCase {
	return &DirectArticleCreateUseCase{
		articleRepository:                  articleRepository,
		studyFieldRepository:               studyFieldRepository,
		authorRepository:                   authorRepository,
		languageRepository:                 languageRepository,
		tagRepository:                      tagRepository,
		articleAuthorAffiliationRepository: articleAuthorAffiliationRepository,
		permissionUseCase:                  permissionUseCase,
		publishRoiSenderTask:               publishRoiSenderTask,
	}
}

func (this *DirectArticleCreateUseCase) Execute(article *entity.ArticlePublishEntity, user *entity.UserBasicEntity, journalID uint) error {
	allowed, err := this.permissionUseCase.Execute(user.Roles, journalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.NewFailResponse(403, "You do not have permission to create or update Article for this journal.")
	}

	if err := validation.ValidateIDs("StudyField", article.StudyFieldsIDs, this.studyFieldRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("CoAuthor", article.CoAuthorsIDs, this.authorRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("Language", []uint{article.LanguageID}, this.languageRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("ArticleAuthorAffiliation", article.ArticleAuthorAffiliationIDs, this.articleAuthorAffiliationRepository.ExistingIds); err != nil {
		return err
	}

	if article.PublishedDate == nil {
		now := time.Now()
		article.PublishedDate = &now
	}

	articleID, err := this.articleRepository.CreatePublishArticle(article, journalID)
	if err != nil {
		return err
	}

	if err := this.publishRoiSenderTask.Run(tasks.PublishRoiSenderPayload{ArticleID: articleID}); err != nil {
		return err
	}

	return nil
}
