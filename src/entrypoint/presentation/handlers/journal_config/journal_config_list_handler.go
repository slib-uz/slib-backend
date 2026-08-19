package journal_config

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalconfigusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_config/schema"
)

type JournalConfigListHandler struct {
	us *JournalConfigUseCase
}

// @inject
func NewJournalConfigListHandler(us *JournalConfigUseCase) *JournalConfigListHandler {
	return &JournalConfigListHandler{us: us}
}

// Handle Get Journal Config List
// @Tags journal-config
// @Accept: application/json
// @Produce: application/json
// @Param creator_id query int false "Creator ID filter"
// @Param journal_id query int false "Journal ID filter"
// @Param is_active query bool false "Is active filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /journal-config/list [get]
// @Security BearerAuth
func (this *JournalConfigListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	page, pageSize := context.GetPagingParams(c)
	creatorID := context.GetIntQueryParam(c, "creator_id", 0)
	journalID := context.GetIntQueryParam(c, "journal_id", 0)
	isActive := context.GetBoolQueryParam(c, "is_active")

	paging, err := this.us.List(uint(creatorID), uint(journalID), isActive, page, pageSize)
	if err != nil {
		return err
	}

	items := make([]*schema.JournalConfigResponse, len(paging.Items))
	for i, item := range paging.Items {
		items[i] = schema.FromJournalConfigEntity(item)
	}

	response := &entity.PagingEntity[schema.JournalConfigResponse]{
		Page:  paging.Page,
		Size:  paging.Size,
		Total: paging.Total,
		Items: items,
	}

	return c.JsonResponse(200, response)
}
