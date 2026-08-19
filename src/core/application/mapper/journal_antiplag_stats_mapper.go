package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func JournalAntiPlagStatsMapper(stats []*entity.JournalAntiPlagStatsEntity) []*entity.JournalAntiPlagStatsWithNameEntity {
	var result = make([]*entity.JournalAntiPlagStatsWithNameEntity, len(stats))
	for i, stat := range stats {
		result[i] = entity.NewJournalAntiPlagStatsWithNameEntity(stat.JournalID, stat.GetJournalName(), stat.Success, stat.Failed)
	}
	return result
}
