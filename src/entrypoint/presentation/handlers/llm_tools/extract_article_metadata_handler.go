package llm_tools

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/llmtoolsusecases"
	appcontext "slib.uz/src/entrypoint/presentation/app/context"
)

type ExtractArticleMetadataHandler struct {
	uc *llmtoolsusecases.ExtractArticleMetadataUseCase
}

// @inject
func NewExtractArticleMetadataHandler(uc *llmtoolsusecases.ExtractArticleMetadataUseCase) *ExtractArticleMetadataHandler {
	return &ExtractArticleMetadataHandler{uc: uc}
}

// Handle
// @Tags llm-tools
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Article PDF"
// @Param journal_id formData int true "Journal ID"
// @Param article_language formData string true "Article language code (uz, ru, en, de, ...)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /llm-tools/article-metadata [post]
func (this *ExtractArticleMetadataHandler) Handle(ctx echo.Context) error {
	c := ctx.(*appcontext.Context)

	file, err := c.FormFile("file")
	if err != nil {
		return response.InvalidFileError
	}

	journalID, err := strconv.ParseUint(strings.TrimSpace(c.FormValue("journal_id")), 10, 64)
	if err != nil {
		return response.InvalidArgument
	}
	articleLanguage := strings.TrimSpace(c.FormValue("article_language"))

	result, err := this.uc.Execute(file, uint(journalID), articleLanguage)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, result)
}
