package roiusecase

import (
	"fmt"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/messaging"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/infrastructure/config"
)

type ROISyncUseCase struct {
	applicationRepository repository.ApplicationRepository
	reviewStageRepo       repository.ReviewStageRepository
	articleRepository     repository.ArticleRepository
	storage               storage.FileStorage
	producer              messaging.Producer
	gateway               gateway.ROIGateway
	config                *config.Config
}

// @inject
func NewROISyncUseCase(
	applicationRepository repository.ApplicationRepository,
	reviewStageRepo repository.ReviewStageRepository,
	articleRepository repository.ArticleRepository,
	storage storage.FileStorage,
	producer messaging.Producer,
	gateway gateway.ROIGateway,
	config *config.Config,
) *ROISyncUseCase {
	return &ROISyncUseCase{
		applicationRepository: applicationRepository,
		reviewStageRepo:       reviewStageRepo,
		articleRepository:     articleRepository,
		storage:               storage,
		producer:              producer,
		gateway:               gateway,
		config:                config,
	}
}

func (this *ROISyncUseCase) Execute(appID uint) (*entity.ROIDataEntity, error) {
	app, err := this.applicationRepository.FindByIDWithRelations(appID)

	if err != nil {
		return nil, err
	}
	// external id for roi
	sourceLink := fmt.Sprintf("%s/article/%d", this.config.FrontendURL, app.ArticleID)

	roi, _err := this.gateway.Publish(app.Article, sourceLink)
	if _err != nil {
		return nil, response.NewOptionalResponse(200, response.RoiUnexpectedError, nil, "")
	}

	if roi.RoiCode == "" {
		return nil, response.NewOptionalResponse(200, response.RoiUnexpectedError, nil, "")
	}

	if roiErr := this.articleRepository.UpdateROI(roi.ID, roi.ExternalID, &roi.RoiCode); roiErr != nil {
		return nil, roiErr
	}

	return roi, nil
}
