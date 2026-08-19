package uzsci

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/uzsciusecases"
	"slib.uz/src/core/domain/entity"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/uzsci/schema"
)

type UzSciCreateApplicationHandler struct {
	uc *uzsciusecases.UzSciCreateApplicationUseCase
}

// @inject
func NewUzSciCreateApplicationHandler(uc *uzsciusecases.UzSciCreateApplicationUseCase) *UzSciCreateApplicationHandler {
	return &UzSciCreateApplicationHandler{uc: uc}
}

// Handle
// @Tags         UzSci
// @Accept       json
// @Produce      json
// @Param        period_id  path  int  true  "Rating period ID"
// @Param        request    body  schema.CreateApplicationRequest  true  "Create application request"
// @Success      201  {object}  response.Response
// @Router       /uzsci/rating-periods/{period_id}/applications [post]
func (this *UzSciCreateApplicationHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	periodID, err := context2.GetIntPathParam(c, "period_id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[schema.CreateApplicationRequest](c)
	if err != nil {
		return err
	}

	answers := make([]entity.UzSciApplicationAnswerEntity, len(body.Answers))
	for i, answer := range body.Answers {
		answers[i] = entity.UzSciApplicationAnswerEntity{
			FormID: answer.FormID,
			Value:  answer.Value,
		}
	}

	if err := this.uc.Execute(uint(periodID), body.JournalID, answers); err != nil {
		return err
	}

	return c.JsonResponse(201, nil)
}
