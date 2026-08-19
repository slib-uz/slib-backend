package middlewares

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JwAuthMiddleware struct {
	authService           *service.UserAuthTokenService
	userProfileRepository repository.UserProfileRepository
}

// @inject
func NewJwAuthMiddleware(authService *service.UserAuthTokenService, userProfileRepository repository.UserProfileRepository) *JwAuthMiddleware {
	return &JwAuthMiddleware{authService: authService, userProfileRepository: userProfileRepository}
}

func (this *JwAuthMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.(*context.Context)

		if err := this.userAuth(c); err != nil {
			return err
		}

		return next(c)
	}
}

func (this *JwAuthMiddleware) userAuth(c *context.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

	if authHeader == "" || !this.isBearer(authHeader) {
		return nil
	}
	tokenStr := this.getBearerToken(authHeader)

	// Access yo'lida grace yo'q: bekor qilingan token darhol rad etiladi.
	userDto, token, err := this.authService.VerifyToken(c.Request().Context(), tokenStr, enum.TokenTypeAccess, 0)
	if err != nil {
		return err
	}

	c.User = userDto
	c.TokenID = token.ID
	c.TokenExp = token.Exp

	this.updateLastOnlineAsync(userDto.ID)
	return nil
}

func (this *JwAuthMiddleware) updateLastOnlineAsync(userID uint) {
	go func() {
		defer func() { recover() }()
		if err := this.userProfileRepository.UpdateLastOnlineAt(userID); err != nil {
			log.Println(err)
		}
	}()
}

func (this *JwAuthMiddleware) isBearer(authHeader string) bool {
	return strings.HasPrefix(authHeader, "Bearer ")
}

func (this *JwAuthMiddleware) getBearerToken(authHeader string) string {
	return strings.TrimPrefix(authHeader, "Bearer ")
}
