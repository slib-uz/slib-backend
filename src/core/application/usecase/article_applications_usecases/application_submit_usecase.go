package usecase

import (
	"log"
	"time"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/validation"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/core/utils"
)

type ApplicationSubmitUseCase struct {
	articleRepository                  repository.ApplicationRepository
	studyFieldRepository               repository.StudyFieldRepository
	authorRepository                   repository.AuthorRepository
	languageRepository                 repository.LanguageRepository
	tagRepository                      repository.TagRepository
	storage                            storage.FileStorage
	sendNotificationService            *service.SendNotificationService
	roleRepo                           repository.UserRoleRepository
	envService                         conf.ConfigAdapter
	ArticleSubmissionUseCase           *ArticleSubmissionUseCase
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository
	ArticleSubmissionRepository        repository.ArticleSubmissionRepository
	journalRepository                  repository.JournalRepository
}

// @inject
func NewApplicationSubmitUseCase(
	articleRepository repository.ApplicationRepository,
	studyFieldRepository repository.StudyFieldRepository,
	authorRepository repository.AuthorRepository,
	languageRepository repository.LanguageRepository,
	tagRepository repository.TagRepository,
	storage storage.FileStorage,
	sendNotificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	envService conf.ConfigAdapter,
	ArticleSubmissionUseCase *ArticleSubmissionUseCase,
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository,
	ArticleSubmissionRepository repository.ArticleSubmissionRepository,
	journalRepository repository.JournalRepository,
) *ApplicationSubmitUseCase {
	return &ApplicationSubmitUseCase{
		articleRepository:                  articleRepository,
		studyFieldRepository:               studyFieldRepository,
		authorRepository:                   authorRepository,
		languageRepository:                 languageRepository,
		tagRepository:                      tagRepository,
		storage:                            storage,
		sendNotificationService:            sendNotificationService,
		roleRepo:                           roleRepo,
		envService:                         envService,
		ArticleSubmissionUseCase:           ArticleSubmissionUseCase,
		articleAuthorAffiliationRepository: articleAuthorAffiliationRepository,
		ArticleSubmissionRepository:        ArticleSubmissionRepository,
		journalRepository:                  journalRepository,
	}
}

func (this *ApplicationSubmitUseCase) Execute(data *entity2.ArticlePublishEntity, journalId uint, userId uint) error {
	if err := validation.ValidateIDs("Journal", []uint{journalId}, this.journalRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("ArticleAuthorAffiliation", data.ArticleAuthorAffiliationIDs, this.articleAuthorAffiliationRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("CoAuthor", data.CoAuthorsIDs, this.authorRepository.ExistingIds); err != nil {
		return err
	}

	authorAffiliationEntity, err := this.ArticleSubmissionUseCase.Execute(data.CoAuthorsIDs, data.ArticleAuthorAffiliationIDs)
	if err != nil {
		return err
	}

	if err := validation.ValidateIDs("StudyField", data.StudyFieldsIDs, this.studyFieldRepository.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("Language", []uint{data.LanguageID}, this.languageRepository.ExistingIds); err != nil {
		return err
	}

	tagIds, err := this.tagRepository.GetOrCreateTags(data.Tags)

	if err != nil {
		return err
	}

	article := mapper.ArticlePublishEntityToCreate(data, journalId, tagIds)

	if err := this.articleRepository.CheckUniqueConstraints(article); err != nil {
		return response.NewFailResponse(400, err.Error())
	}

	app, err := this.ArticleSubmissionRepository.Create(article, authorAffiliationEntity, data.ArticleAuthorAffiliationIDs, userId, utils.MakeDeadline(time.Now(), this.envService.GetReviewDeadlineDays()))
	if err != nil {
		return err
	}

	this.sendNotification(app)

	return nil
}

func (this *ApplicationSubmitUseCase) sendNotification(application *entity2.ApplicationEntity) {
	journalMemberIds, err := this.roleRepo.GetJournalMemberIds(application.JournalID)
	if err != nil {
		log.Printf("Failed to get journal member IDs: %v\n", err)
		return
	}

	if err := this.sendNotificationService.NewApplication(journalMemberIds, application.ID, application.Number); err != nil {
		log.Printf("Failed to send notification: %v\n", err)
	}
}
