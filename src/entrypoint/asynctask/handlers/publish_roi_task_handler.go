package handlers

import (
	"context"
	"fmt"
	"log"

	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type PublishRoiTaskHandler struct {
	publishedArticleRepo repository.PublishedArticleRepository
	articleRepo          repository.ArticleRepository
	roiGateway           gateway.ROIGateway
	frontendURL          string
}

// @inject
func NewPublishRoiTaskHandler(
	publishedArticleRepo repository.PublishedArticleRepository,
	articleRepo repository.ArticleRepository,
	roiGateway gateway.ROIGateway,
	configAdapter conf.ConfigAdapter,
) *PublishRoiTaskHandler {
	return &PublishRoiTaskHandler{
		publishedArticleRepo: publishedArticleRepo,
		articleRepo:          articleRepo,
		roiGateway:           roiGateway,
		frontendURL:          configAdapter.GetFrontendURL(),
	}
}

func (h *PublishRoiTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	payload, err := UnmarshalPayload[tasks.PublishRoiSenderPayload](task)
	if err != nil {
		log.Println("[PublishRoiTaskHandler] UnmarshalPayload:", err.Error())
		return err
	}

	article, err := h.publishedArticleRepo.GetByIDWithRelations(payload.ArticleID)
	if err != nil {
		log.Println("[PublishRoiTaskHandler] GetByIDWithRelations:", err.Error())
		return err
	}

	sourceLink := fmt.Sprintf("%s/article/%d", h.frontendURL, article.ID)

	roi, err := h.roiGateway.Publish(article, sourceLink)
	if err != nil {
		log.Println("[PublishRoiTaskHandler] Publish:", err.Error())
		return err
	}

	if roi.RoiCode == "" {
		log.Println("[PublishRoiTaskHandler] empty ROI code for article:", payload.ArticleID)
		return nil
	}

	if err := h.articleRepo.UpdateROI(roi.ID, roi.ExternalID, &roi.RoiCode); err != nil {
		log.Println("[PublishRoiTaskHandler] UpdateROI:", err.Error())
		return err
	}

	return nil
}
