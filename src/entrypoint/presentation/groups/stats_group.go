package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/stats"
)

type StatsGroup struct {
	purchaseStatsByJournalHandler        *stats.PurchaseStatsByJournalHandler
	antiPlagStatsByJournalHandler        *stats.AntiPlagStatsByJournalHandler
	spellcheckStatsByJournalHandler      *stats.SpellcheckStatsByJournalHandler
	reviewStageStatisticsHandler         *stats.ReviewStageStatisticsHandler
	reviewStageOverdueStatisticsHandler  *stats.ReviewStageOverdueStatisticsHandler
	reviewStageOverdueListHandler        *stats.ReviewStageOverdueListHandler
}

// @inject
func NewStatsGroup(
	purchaseStatsByJournalHandler *stats.PurchaseStatsByJournalHandler,
	antiPlagStatsByJournalHandler *stats.AntiPlagStatsByJournalHandler,
	spellcheckStatsByJournalHandler *stats.SpellcheckStatsByJournalHandler,
	reviewStageStatisticsHandler *stats.ReviewStageStatisticsHandler,
	reviewStageOverdueStatisticsHandler *stats.ReviewStageOverdueStatisticsHandler,
	reviewStageOverdueListHandler *stats.ReviewStageOverdueListHandler,
) *StatsGroup {
	return &StatsGroup{
		purchaseStatsByJournalHandler:       purchaseStatsByJournalHandler,
		antiPlagStatsByJournalHandler:       antiPlagStatsByJournalHandler,
		spellcheckStatsByJournalHandler:     spellcheckStatsByJournalHandler,
		reviewStageStatisticsHandler:        reviewStageStatisticsHandler,
		reviewStageOverdueStatisticsHandler: reviewStageOverdueStatisticsHandler,
		reviewStageOverdueListHandler:       reviewStageOverdueListHandler,
	}
}

func (this *StatsGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/purchase", this.purchaseStatsByJournalHandler.Handle)
	group.GET("/antiplag", this.antiPlagStatsByJournalHandler.Handle)
	group.GET("/spellcheck", this.spellcheckStatsByJournalHandler.Handle)
	group.GET("/review-stage", this.reviewStageStatisticsHandler.Handle)
	group.GET("/review-stage/overdue", this.reviewStageOverdueStatisticsHandler.Handle)
	group.GET("/review-stage/overdue/list", this.reviewStageOverdueListHandler.Handle)
}
