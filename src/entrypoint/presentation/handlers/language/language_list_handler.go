package language

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/languageusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type LanguageListHandler struct {
	uc *LanguageListUseCase
}

// @inject
func NewLanguageListHandler(uc *LanguageListUseCase) *LanguageListHandler {
	return &LanguageListHandler{uc: uc}
}

// Handle
// @Tags language
// @Accept json
// @Produce json
// @Success 200 {array} entity.LanguageEntity
// @Router /language/list [get]
func (this *LanguageListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	languages, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, languages)

}
