package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func JournalSpellcheckStatsMapper(stats []*entity.JournalSpellcheckStatsEntity) []*entity.JournalSpellcheckStatsWithNameEntity {
	var result = make([]*entity.JournalSpellcheckStatsWithNameEntity, len(stats))
	for i, stat := range stats {
		result[i] = entity.NewJournalSpellcheckStatsWithNameEntity(stat.JournalID, stat.GetJournalName(), stat.Success, stat.Failed)
	}
	return result
}
