package news

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/newsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type NewsRetrieveHandler struct {
	uc *NewsRetrieveUseCase
}

// @inject
func NewNewsRetrieveHandler(uc *NewsRetrieveUseCase) *NewsRetrieveHandler {
	return &NewsRetrieveHandler{uc: uc}
}

// Handle NewsRetrieveHandler
// @Tags news
// @Accept json
// @Produce json
// @Param X-Anonymous-Token header string false "Anon JWT token"
// @Param newsId path int true "News ID"
// @Success 200 {object} entity.NewsEntity
// @Failure 400 {object} response.Response
// @Router /news/retrieve/{newsId} [get]
func (this *NewsRetrieveHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	userKey, err := this.userKey(c)
	if err != nil {
		return err
	}

	newsId, err := context2.GetIntPathParam(c, "newsId")
	if err != nil {
		return err
	}

	news, err := this.uc.Execute(uint(newsId), userKey)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, news)
}

func (this *NewsRetrieveHandler) userKey(c *context2.Context) (string, error) {
	if c.User == nil {
		if c.AnonymousID == "" {
			return "", response.NewFailResponse(400, "Anonymous token is required")
		}
		return c.AnonymousID, nil
	}
	return strconv.Itoa(int(c.User.ID)), nil
}
