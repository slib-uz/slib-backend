package handlers

import (
	"context"
	"fmt"
	"log"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type UpdateRoiTaskHandler struct {
	publishedArticleRepo repository.PublishedArticleRepository
	roiGateway           gateway.ROIGateway
	frontendURL          string
}

// @inject
func NewUpdateRoiTaskHandler(
	publishedArticleRepo repository.PublishedArticleRepository,
	roiGateway gateway.ROIGateway,
	configAdapter conf.ConfigAdapter,
) *UpdateRoiTaskHandler {
	return &UpdateRoiTaskHandler{
		publishedArticleRepo: publishedArticleRepo,
		roiGateway:           roiGateway,
		frontendURL:          configAdapter.GetFrontendURL(),
	}
}

func (h *UpdateRoiTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	payload, err := UnmarshalPayload[tasks.UpdateRoiSenderPayload](task)
	if err != nil {
		log.Println("[UpdateRoiTaskHandler] UnmarshalPayload:", err.Error())
		return err
	}

	article, err := h.publishedArticleRepo.GetByIDWithRelations(payload.ArticleID)
	if err != nil {
		log.Println("[UpdateRoiTaskHandler] GetByIDWithRelations:", err.Error())
		return err
	}

	updateArticle := mapper.ArticleEntityToROIArticleUpdateEntity(article, fmt.Sprintf("%s/article/%d", h.frontendURL, article.ID))
	if updateArticle.ROI == "" {
		log.Println("[UpdateRoiTaskHandler] Invalid ROI for article:", payload.ArticleID)
		return nil
	}

	if err := h.roiGateway.UpdateArticle(updateArticle); err != nil {
		log.Println("[UpdateRoiTaskHandler] UpdateArticle:", err.Error())
		return err
	}

	return nil
}
