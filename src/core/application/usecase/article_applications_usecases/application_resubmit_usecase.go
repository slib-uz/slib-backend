package usecase

import (
	"errors"
	"fmt"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/validation"
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type ApplicationResubmitUseCase struct {
	articleRepository                  repository.ArticleRepository
	authorRepository                   repository.AuthorRepository
	studyFieldRepository               repository.StudyFieldRepository
	languageRepository                 repository.LanguageRepository
	tagRepository                      repository.TagRepository
	storage                            storage.FileStorage
	applicationRepository              repository.ApplicationRepository
	reviewStageRepository              repository.ReviewStageRepository
	sendNotificationService            *service.SendNotificationService
	roleRepo                           repository.UserRoleRepository
	userRepository                     repository.UserRepository
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository
}

// @inject
func NewApplicationResubmitUseCase(
	articleRepository repository.ArticleRepository,
	authorRepository repository.AuthorRepository,
	studyFieldRepository repository.StudyFieldRepository,
	languageRepository repository.LanguageRepository,
	tagRepository repository.TagRepository,
	storage storage.FileStorage,
	applicationRepository repository.ApplicationRepository,
	reviewStageRepository repository.ReviewStageRepository,
	sendNotificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	userRepository repository.UserRepository,
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository,
) *ApplicationResubmitUseCase {
	return &ApplicationResubmitUseCase{
		articleRepository:                  articleRepository,
		authorRepository:                   authorRepository,
		studyFieldRepository:               studyFieldRepository,
		languageRepository:                 languageRepository,
		tagRepository:                      tagRepository,
		storage:                            storage,
		applicationRepository:              applicationRepository,
		reviewStageRepository:              reviewStageRepository,
		sendNotificationService:            sendNotificationService,
		roleRepo:                           roleRepo,
		userRepository:                     userRepository,
		articleAuthorAffiliationRepository: articleAuthorAffiliationRepository,
	}
}

func (this *ApplicationResubmitUseCase) Execute(applicationID uint, article *entity2.ArticleUpdateEntity, userId uint) error {

	if article.CoAuthorsIDs != nil {
		if err := validation.ValidateIDs("CoAuthor", *article.CoAuthorsIDs, this.authorRepository.ExistingIds); err != nil {
			return err
		}
		if err := this.checkCoAuthorsIds(*article.CoAuthorsIDs, userId); err != nil {
			return err
		}
	}

	if article.AffiliationsIDs != nil {
		if err := validation.ValidateIDs("ArticleAuthorAffiliation", *article.AffiliationsIDs, this.articleAuthorAffiliationRepository.ExistingIds); err != nil {
			return err
		}
	}

	if article.StudyFieldsIDs != nil {
		if err := validation.ValidateIDs("StudyField", *article.StudyFieldsIDs, this.studyFieldRepository.ExistingIds); err != nil {
			return err
		}
	}

	if article.LanguageID != nil {
		if err := validation.ValidateIDs("Language", []uint{*article.LanguageID}, this.languageRepository.ExistingIds); err != nil {
			return err
		}
	}

	if err := this.checkReviewState(applicationID); err != nil {
		return err
	}

	if err := this.checkUserAccess(applicationID, userId); err != nil {
		return err
	}

	if err := this.applicationRepository.Update(applicationID, article); err != nil {
		return err
	}

	app, _ := this.applicationRepository.FindByID(applicationID)
	this.sendNotification(app)

	return nil
}

func (this *ApplicationResubmitUseCase) sendNotification(application *entity2.ApplicationEntity) {
	journalMemberIds, err := this.roleRepo.GetJournalMemberIds(application.JournalID)
	if err != nil {
		fmt.Println("Failed to get journal members(ApplicationResubmitUseCase)", err)
		return
	}

	if err := this.sendNotificationService.ReSent(journalMemberIds, application.ID, application.Number); err != nil {
		fmt.Println("Failed to send notification(ApplicationResubmitUseCase)", err)
	}
}

func (this *ApplicationResubmitUseCase) checkUserAccess(applicationID, userId uint) error {
	application, err := this.applicationRepository.FindByID(applicationID)
	if err != nil {
		return err
	}
	if application.UserID != userId {
		return response.NewFailResponse(403, "You do not have permission to resubmit this application")
	}
	return nil
}

func (this *ApplicationResubmitUseCase) checkReviewState(applicationID uint) error {
	_, err := this.reviewStageRepository.GetBy(applicationID, enum2.TechnicalReviewStage, enum2.StatusRejected)
	if errors.Is(err, response.NotFoundError) {
		return response.NewFailResponse(400, "Application is not rejected in technical review stage")
	}
	return err
}

func (this *ApplicationResubmitUseCase) checkCoAuthorsIds(ids []uint, userId uint) error {
	user, err := this.userRepository.GetById(userId)

	if err != nil {
		return err
	}

	ownAuthorId, err := this.authorRepository.GetOwnCoAuthorId(user.ScienceID)
	for _, id := range ids {
		if id == ownAuthorId {
			return nil
		}
	}
	return response.NewFailResponse(400, "Applicant AuthorID must be in CoAuthorsIds field")
}
