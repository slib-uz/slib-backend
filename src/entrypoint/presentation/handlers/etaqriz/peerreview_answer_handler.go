package etaqriz

import (
	"fmt"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/etaqrizusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/etaqriz/schema"
)

type PeerReviewAnswerHandler struct {
	uc *ChangeStatusUseCase
}

// @inject
func NewPeerReviewAnswerHandler(uc *ChangeStatusUseCase) *PeerReviewAnswerHandler {
	return &PeerReviewAnswerHandler{uc: uc}
}

// Handle
// @Tags e-taqriz
// @Accept json
// @Produce json
// @Router /etaqriz/callback [post]
func (this *PeerReviewAnswerHandler) Handle(ctx echo.Context) error {

	data, err := context.GetBody[schema.PeerReviewAnswerSchema](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data.SubmissionID, data.Status, data.Deadline, data.OldDeadline, data.Version); err != nil {
		fmt.Println("Error executing use case:", err)
		return err
	}

	return ctx.JSON(200, data)

}
