package handlers

import (
	"context"

	usecase "slib.uz/src/core/application/usecase/article_applications_usecases"
	"slib.uz/src/core/domain/entity"
)

type SendArticleKafkaTaskHandler struct {
	uc *usecase.ApplicationArticleInvokeUseCase
}

// @inject
func NewSendArticleKafkaTaskHandler(
	uc *usecase.ApplicationArticleInvokeUseCase,
) *SendArticleKafkaTaskHandler {
	return &SendArticleKafkaTaskHandler{uc: uc}
}

func (h *SendArticleKafkaTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	appID, err := UnmarshalPayload[uint](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(appID)
}
