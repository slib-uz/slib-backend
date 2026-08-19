package spellcheck

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/spellcheck/schema"
)

type SpellCheckHandler struct {
	uc *SpellCheckUsecase
}

// @inject
func NewSpellCheckHandler(uc *SpellCheckUsecase) *SpellCheckHandler {
	return &SpellCheckHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.SpellCheckRequest true "Spell check request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /article-application/spellcheck [post]
func (this *SpellCheckHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.SpellCheckRequest](ctx)
	if err != nil {
		return err
	}

	id, err := this.uc.Execute(c.User.ID, data.ReviewStageID, data.ApplicationID)

	if err != nil {
		return err
	}
	return c.JsonResponse(200, map[string]any{"id": id})
}
