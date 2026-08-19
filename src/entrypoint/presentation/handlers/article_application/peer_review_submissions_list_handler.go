package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/peerreviewusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type PeerReviewSubmissionsListHandler struct {
	uc *PeerReviewSubmissionsListUseCase
}

// @inject
func NewPeerReviewSubmissionsListHandler(uc *PeerReviewSubmissionsListUseCase) *PeerReviewSubmissionsListHandler {
	return &PeerReviewSubmissionsListHandler{uc: uc}
}

// Handle godoc
// @Summary      Get peer review submissions list
// @Description  Get peer review submissions list by application ID
// @Tags         article-application
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        applicationId path int true "Application ID"
// @Success      200 {array} entity.PeerReviewSubmissionEntity
// @Router       /article-application/{applicationId}/peer-review/submissions/list [get]
func (this *PeerReviewSubmissionsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	applicationId, err := context2.GetIntPathParam(c, "applicationId")
	if err != nil {
		return err
	}

	submissions, err := this.uc.Execute(uint(applicationId))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, submissions)
}
