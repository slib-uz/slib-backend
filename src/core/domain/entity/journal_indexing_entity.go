package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type JournalIndexingEntity struct {
	ID           uint              `json:"id"`
	IndexingType enum.IndexingType `json:"indexing"`
	URL          string            `json:"url"`
	JournalID    uint              `json:"journal_id"`
	Journal      *JournalEntity    `json:"journal,omitempty"`
}

func NewJournalIndexingEntity(id uint, indexingType enum.IndexingType, URL string, journalID uint, journal *JournalEntity) *JournalIndexingEntity {
	return &JournalIndexingEntity{ID: id, IndexingType: indexingType, URL: URL, JournalID: journalID, Journal: journal}
}

func NewJournalIndexingCreateEntity(indexingType enum.IndexingType, url string) *JournalIndexingEntity {
	return &JournalIndexingEntity{IndexingType: indexingType, URL: url}
}
