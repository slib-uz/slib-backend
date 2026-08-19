package articleusecases

import (
	"fmt"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authorusecases"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/core/application/validation"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// IntegrationArticleCreateUseCase — tashqi tizimlardan (basic auth orqali)
// article yaratish uchun usecase. JWT autentifikatsiyasi va
// journal membership tekshiruvi talab qilinmaydi.
type IntegrationArticleCreateUseCase struct {
	articleRepository                  repository.ArticleRepository
	studyFieldRepository               repository.StudyFieldRepository
	authorRepository                   repository.AuthorRepository
	languageRepository                 repository.LanguageRepository
	tagRepository                      repository.TagRepository
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository
	journalRepository                  repository.JournalRepository
	uploadBase64FileUc                 *uploadusecases.UploadBase64FileUseCase
	authorByScienceIDUc                *authorusecases.AuthorByScienceIDUseCase
}

// @inject
func NewIntegrationArticleCreateUseCase(
	articleRepository repository.ArticleRepository,
	studyFieldRepository repository.StudyFieldRepository,
	authorRepository repository.AuthorRepository,
	languageRepository repository.LanguageRepository,
	tagRepository repository.TagRepository,
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository,
	journalRepository repository.JournalRepository,
	uploadBase64FileUc *uploadusecases.UploadBase64FileUseCase,
	authorByScienceIDUc *authorusecases.AuthorByScienceIDUseCase,
) *IntegrationArticleCreateUseCase {
	return &IntegrationArticleCreateUseCase{
		articleRepository:                  articleRepository,
		studyFieldRepository:               studyFieldRepository,
		authorRepository:                   authorRepository,
		languageRepository:                 languageRepository,
		tagRepository:                      tagRepository,
		articleAuthorAffiliationRepository: articleAuthorAffiliationRepository,
		journalRepository:                  journalRepository,
		uploadBase64FileUc:                 uploadBase64FileUc,
		authorByScienceIDUc:                authorByScienceIDUc,
	}
}

func (this *IntegrationArticleCreateUseCase) Execute(article *entity.ArticlePublishEntity, journalID uint) error {
	if err := validation.ValidateIDs("Journal", []uint{journalID}, this.journalRepository.ExistingIds); err != nil {
		return err
	}

	// Unique constraint tekshiruvlari
	if article.DOI != nil && *article.DOI != "" {
		exists, err := this.articleRepository.CheckDOIExists(*article.DOI)
		if err != nil {
			return err
		}
		if exists {
			return response.NewFailResponse(409, fmt.Sprintf("article with DOI %s already exists", *article.DOI))
		}
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

	for i, affiliation := range article.AuthorAffiliations {
		if affiliation.ScienceID == "" {
			return response.NewFailResponse(400, fmt.Sprintf("affiliations[%d]: science_id is required", i))
		}
		if affiliation.OrganizationTin == "" {
			return response.NewFailResponse(400, fmt.Sprintf("affiliations[%d]: organization_tin is required", i))
		}
		if affiliation.OrganizationName == "" {
			return response.NewFailResponse(400, fmt.Sprintf("affiliations[%d]: organization_name is required", i))
		}

		author, err := this.authorByScienceIDUc.Execute(affiliation.ScienceID)
		if err != nil {
			return err
		}
		affiliation.AuthorID = author.ID
	}

	// Validatsiyadan o'tgach, fayllarni yuklash
	contentFilePath, err := this.uploadBase64FileUc.Execute(article.ContentFile, enum.FolderArticle, enum.BucketPrivate)
	if err != nil {
		return err
	}
	article.ContentFile = contentFilePath

	var expertConclusionFilePath *string
	if article.ExpertConclusionFile != nil && *article.ExpertConclusionFile != "" {
		path, err := this.uploadBase64FileUc.Execute(*article.ExpertConclusionFile, enum.FolderArticleExpertConclusion, enum.BucketPublic)
		if err != nil {
			return err
		}
		expertConclusionFilePath = &path
	}
	article.ExpertConclusionFile = expertConclusionFilePath

	if article.PublishedDate == nil {
		now := time.Now()
		article.PublishedDate = &now
	}

	_, err = this.articleRepository.CreatePublishArticleWithAffiliations(article, journalID)
	if err != nil {
		return err
	}

	return nil
}
