package roiusecase

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type ArticlePublishROIUseCase struct {
	roi             *ROISyncUseCase
	roiSenderTask   *tasks.SetRoiSenderTask
	applicationRepo repository.ApplicationRepository
	storage         storage.FileStorage
	permission      *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewArticlePushROIUseCase(
	roi *ROISyncUseCase,
	roiSenderTask *tasks.SetRoiSenderTask,
	applicationRepo repository.ApplicationRepository,
	storage storage.FileStorage,
	permission *permissionusecases.JournalMemberPermissionUseCase,
) *ArticlePublishROIUseCase {
	return &ArticlePublishROIUseCase{
		roi:             roi,
		roiSenderTask:   roiSenderTask,
		applicationRepo: applicationRepo,
		storage:         storage,
		permission:      permission,
	}
}

func (this *ArticlePublishROIUseCase) Execute(user *entity.UserBasicEntity, applicationID uint, finalFilePath string) (*entity.ROIDataEntity, error) {
	application, err := this.applicationRepo.GetWithArticleAndJournal(applicationID)
	if err != nil {
		return nil, err
	}
	allowed, err := this.permission.Execute(user.Roles, application.Article.JournalID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, response.PermissionDeniedError
	}

	if err := this.applicationRepo.Publish(application.Article.ID, applicationID, finalFilePath); err != nil {
		return nil, err
	}

	result, err := this.roi.Execute(applicationID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
