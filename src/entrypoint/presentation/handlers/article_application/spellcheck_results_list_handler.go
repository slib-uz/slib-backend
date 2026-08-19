package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/spellcheckusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type SpellCheckResultsListHandler struct {
	uc *SpellCheckResultsListUsecase
}

// @inject
func NewSpellCheckResultsListHandler(uc *SpellCheckResultsListUsecase) *SpellCheckResultsListHandler {
	return &SpellCheckResultsListHandler{uc: uc}
}

// Handle
// @Tags spellcheck
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param applicationId path int true "Application ID"
// @Param reviewStageId query int false "Spell Check ID"
// @Success 200 {array} entity.SpellCheckResultEntity
// @Router /article-application/{applicationId}/spellcheck/results [get]
func (this *SpellCheckResultsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	applicationId, err := context2.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	results, err := this.uc.Execute(
		c.User,
		uint(applicationId),
		uint(context2.GetIntQueryParam(ctx, "reviewStageId", 0)),
	)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, results)
}
