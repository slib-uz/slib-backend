package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/core/domain/entity"
)

type ArticleViewsCountUpdateTaskHandler struct {
	uc *articleusecases.UpdateArticleViewsCountUseCase
}

// @inject
func NewArticleViewsCountUpdateTaskHandler(uc *articleusecases.UpdateArticleViewsCountUseCase) *ArticleViewsCountUpdateTaskHandler {
	return &ArticleViewsCountUpdateTaskHandler{uc: uc}
}

func (h *ArticleViewsCountUpdateTaskHandler) ProcessTask(_ context.Context, _ *entity.AsyncTask) error {
	return h.uc.Execute()
}
