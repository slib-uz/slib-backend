package guide

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/guideusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type GuideRetrieveHandler struct {
	uc *GuideRetrieveUseCase
}

// @inject
func NewGuideRetrieveHandler(uc *GuideRetrieveUseCase) *GuideRetrieveHandler {
	return &GuideRetrieveHandler{uc: uc}
}

// Handle GuideRetrieveHandler
// @Tags guide
// @Accept json
// @Produce json
// @Param guideId path uint true "Guide ID"
// @Success 200 {object} entity.GuideRetrieveEntity
// @Failure 404 {object} response.Response
// @Router /guides/retrieve/{guideId} [get]
func (this *GuideRetrieveHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	guideId, err := context2.GetIntPathParam(c, "guideId")
	if err != nil {
		return err
	}

	guide, err := this.uc.Execute(uint(guideId))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, guide)
}
