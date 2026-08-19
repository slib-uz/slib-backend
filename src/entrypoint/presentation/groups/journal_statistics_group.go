package groups

import (
	"github.com/labstack/echo/v4"
	journalstatistics "slib.uz/src/entrypoint/presentation/handlers/journal_statistics"
)

type JournalStatisticsGroup struct {
	list  *journalstatistics.JournalStatisticsListHandler
	excel *journalstatistics.JournalStatisticsListExcelHandler
	detail *journalstatistics.JournalStatisticsDetailHandler
}

// @inject
func NewJournalStatisticsGroup(
	list *journalstatistics.JournalStatisticsListHandler,
	excel *journalstatistics.JournalStatisticsListExcelHandler,
	detail *journalstatistics.JournalStatisticsDetailHandler,
) *JournalStatisticsGroup {
	return &JournalStatisticsGroup{list: list, excel: excel, detail: detail}
}

func (this *JournalStatisticsGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.list.Handle)
	group.GET("/list/excel", this.excel.Handle)
	group.GET("/detail/:journal_id", this.detail.Handle)
}
