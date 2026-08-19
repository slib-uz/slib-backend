package social

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/socialusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type SocialAllHandler struct {
	usecase *socialusecases.SocialAllUseCase
}

// @inject
func NewSocialAllHandler(usecase *socialusecases.SocialAllUseCase) *SocialAllHandler {
	return &SocialAllHandler{
		usecase: usecase,
	}
}

// Handle SocialAllHandler
// @Tags Social
// @Accept json
// @Produce json
// @Success 200
// @Router /social/all [get]
func (this *SocialAllHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	result, err := this.usecase.Execute()
	if err != nil {
		return c.JsonResponse(404, err.Error())
	}
	return c.JsonResponse(200, result)
}
