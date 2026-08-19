package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/peerreviewusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type PeerReviewSubmissionReviewsHandler struct {
	uc *PeerReviewSubmissionReviewsUseCase
}

// @inject
func NewPeerReviewSubmissionReviewsHandler(uc *PeerReviewSubmissionReviewsUseCase) *PeerReviewSubmissionReviewsHandler {
	return &PeerReviewSubmissionReviewsHandler{uc: uc}
}

// Handle godoc
// @Tags article-application
// @Security BearerAuth
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param external_id path int true "Submission External ID"
// @Success 200 {object} entity.PeerReviewSubmissionEntity
// @Router /article-application/peer-review/reviews/{external_id} [get]
func (this *PeerReviewSubmissionReviewsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	submissionExternalID, err := context2.GetIntPathParam(c, "external_id")
	if err != nil {
		return err
	}

	submission, err := this.uc.Execute(uint(submissionExternalID))
	if err != nil {
		return err
	}
	return c.JSON(200, submission)
}
