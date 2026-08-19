package roiusecase

import (
	"fmt"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

// ROIPublishByArticleUseCase — berilgan article'ni ROI'ga sinxron publish qilib,
// uning roi va external_id'sini yangilaydi (admin uchun).
type ROIPublishByArticleUseCase struct {
	publishedArticleRepo repository.PublishedArticleRepository
	articleRepository    repository.ArticleRepository
	gateway              gateway.ROIGateway
	frontendURL          string
}

// @inject
func NewROIPublishByArticleUseCase(
	publishedArticleRepo repository.PublishedArticleRepository,
	articleRepository repository.ArticleRepository,
	roiGateway gateway.ROIGateway,
	configAdapter conf.ConfigAdapter,
) *ROIPublishByArticleUseCase {
	return &ROIPublishByArticleUseCase{
		publishedArticleRepo: publishedArticleRepo,
		articleRepository:    articleRepository,
		gateway:              roiGateway,
		frontendURL:          configAdapter.GetFrontendURL(),
	}
}

func (this *ROIPublishByArticleUseCase) Execute(articleID uint) (*entity.ROIDataEntity, error) {
	article, err := this.publishedArticleRepo.GetByIDWithRelations(articleID)
	if err != nil {
		return nil, err
	}

	sourceLink := fmt.Sprintf("%s/article/%d", this.frontendURL, article.ID)

	roi, err := this.gateway.Publish(article, sourceLink)
	if err != nil {
		return nil, err
	}

	if roi.RoiCode == "" {
		return nil, response.NewFailResponse(400, "ROI data not found")
	}

	if err := this.articleRepository.UpdateROI(roi.ID, roi.ExternalID, &roi.RoiCode); err != nil {
		return nil, err
	}

	return roi, nil
}
