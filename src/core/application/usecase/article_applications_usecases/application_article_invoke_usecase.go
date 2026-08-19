package usecase

import (
	"log"
	"time"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/messaging"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
)

type ApplicationArticleInvokeUseCase struct {
	applicationRepository repository.ApplicationRepository
	reviewStageRepo       repository.ReviewStageRepository
	articleRepository     repository.ArticleRepository
	storage               storage.FileStorage
	producer              messaging.Producer
	gateway               gateway.ROIGateway
	config                *config.Config
}

// @inject
func NewApplicationArticleInvokeUseCase(applicationRepository repository.ApplicationRepository, reviewStageRepo repository.ReviewStageRepository, articleRepository repository.ArticleRepository, storage storage.FileStorage, producer messaging.Producer, gateway gateway.ROIGateway, config *config.Config) *ApplicationArticleInvokeUseCase {
	return &ApplicationArticleInvokeUseCase{
		applicationRepository: applicationRepository,
		reviewStageRepo:       reviewStageRepo,
		articleRepository:     articleRepository,
		storage:               storage,
		producer:              producer,
		gateway:               gateway,
		config:                config,
	}
}

func (this *ApplicationArticleInvokeUseCase) Execute(appID uint) error {
	app, err := this.applicationRepository.FindByIDWithRelations(appID)

	if err != nil {
		return err
	}

	KEY := "article.slib.uz"
	EXPIRES := 7 * 24 * time.Hour
	file, minioErr := this.storage.PresignedURL(enum.BucketPrivate, app.Article.ContentFile, EXPIRES)
	if minioErr != nil {
		return minioErr
	}

	msgBytes, marshalErr := mapper.ArticleToMessagingMapper(app.Article, file)
	if marshalErr != nil {
		return marshalErr
	}
	if kafkaErr := this.producer.PublishSync(this.config.KafkaTopic, KEY, msgBytes); kafkaErr != nil {
		return kafkaErr
	}
	err = this.articleRepository.UpdateStatus(app.ArticleID, true)
	if err != nil {
		log.Println("[Kafka]:", err)
		return err
	}
	return nil
}
