package middlewares

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ClientJournalAuthMiddleware struct{}

// @inject
func NewClientJournalAuthMiddleware() *ClientJournalAuthMiddleware {
	return &ClientJournalAuthMiddleware{}
}

func (m *ClientJournalAuthMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)

		client := ctx.GetClient()
		if client == nil {

			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		headerJournalID := ctx.Request().Header.Get("JournalID")
		if headerJournalID == "" {
			return ctx.JSON(response.UnauthorizedError.Status, response.JournalIDHeaderRequired)
		}

		parsedJournalID, err := strconv.ParseUint(headerJournalID, 10, 64)
		if err != nil {
			return ctx.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
		}

		journalIDUint := uint(parsedJournalID)

		if client.JournalID == nil || *client.JournalID != journalIDUint {
			return ctx.JSON(response.UnauthorizedError.Status, response.JournalIDHeaderMismatch)
		}

		return next(ctx)
	}
}
