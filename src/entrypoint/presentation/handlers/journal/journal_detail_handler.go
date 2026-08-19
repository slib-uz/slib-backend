package journal

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalDetailHandler struct {
	uc *JournalDetailUseCase
}

// @inject
func NewJournalDetailHandler(uc *JournalDetailUseCase) *JournalDetailHandler {
	return &JournalDetailHandler{uc: uc}
}

// Handle JournalDetailHandler
// @Tags journal
// @Accept json
// @Produce json
// @Param X-Anonymous-Token  header  string  false  "Anon JWT token"
// @Param id path uint true "Journal ID"
// @Success 200 {object} entity.JournalEntity
// @Failure 404 {object} response.Response
// @Router /journal/detail/{id} [get]
func (this *JournalDetailHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context.Context)

	userKey, err := this.userKey(c)
	if err != nil {
		return err
	}

	journal, err := this.uc.Execute(this.GetID(c), userKey)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, journal)

}

func (this *JournalDetailHandler) GetID(ctx echo.Context) uint {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(response.NewFailResponse(400, "Invalid ID"))
	}
	return uint(id)
}

func (this *JournalDetailHandler) userKey(c *context.Context) (string, error) {
	if c.User == nil {
		if c.AnonymousID == "" {
			return "", response.NewFailResponse(400, "Anonymous token is required")
		}
		return c.AnonymousID, nil
	}
	return strconv.Itoa(int(c.User.ID)), nil
}
