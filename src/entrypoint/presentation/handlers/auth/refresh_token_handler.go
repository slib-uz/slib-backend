package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type RefreshTokenHandler struct {
	uc *authusecases.RefreshTokenUseCase
}

// @inject
func NewRefreshTokenHandler(usecase *authusecases.RefreshTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{uc: usecase}
}

// Handle godoc
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data           body   schema.RefreshTokenRequest  false  "Refresh token"
// @Param        refresh_token  query  string                      false  "Refresh token (eskirgan, body'dan foydalaning)"
// @Success      200  {object}  schema.RefreshTokenResponse
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/refresh-token [post]
func (this *RefreshTokenHandler) Handle(ctx echo.Context) error {
	refreshToken := ""
	if req, err := context2.GetBody[schema.RefreshTokenRequest](ctx); err == nil {
		refreshToken = req.RefreshToken
	}
	// Orqaga moslik: 1-bosqichda eski frontend query parametr yuboradi.
	// 2-bosqichda bu tarmoq olib tashlanadi.
	if refreshToken == "" {
		refreshToken = ctx.QueryParam("refresh_token")
	}
	if refreshToken == "" {
		return response.NewFailResponse(400, "refresh_token talab qilinadi")
	}

	result, err := this.uc.Execute(ctx.Request().Context(), refreshToken)
	if err != nil {
		return err
	}

	// Javob ataylab yassi (c.JsonResponse orqali o'ralmaydi): shu endpoint
	// bu shaklni har doim qaytargan va ishlab turgan frontend shunga bog'liq.
	// O'rash mijoz uchun buzilish bo'lardi — aynan bosqichli rollout oldini
	// olmoqchi bo'lgan narsa. Shakl schema.RefreshTokenResponse da mahkamlangan.
	return ctx.JSON(http.StatusOK, map[string]any{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}
