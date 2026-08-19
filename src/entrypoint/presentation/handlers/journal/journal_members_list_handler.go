package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalMembersListHandler struct {
	uc *JournalMembersListUsecase
}

// @inject
func NewJournalMembersListHandler(uc *JournalMembersListUsecase) *JournalMembersListHandler {
	return &JournalMembersListHandler{uc: uc}
}

// Handle
// @Tags journal-manage
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journalId path int true "Journal ID"
// @Success 200 {array} entity.UserRoleWithBasicUserEntity
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /journal-manage/{journalId}/members [get]
func (this *JournalMembersListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalId, err := context2.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	members, err := this.uc.Execute(c.User, uint(journalId))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, members)

}
