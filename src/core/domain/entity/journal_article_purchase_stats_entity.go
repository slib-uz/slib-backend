package entity

import "encoding/json"

type JournalArticlePurchaseStatsEntity struct {
	JournalID   uint   `json:"journal_id"`
	JournalName []byte `json:"journal_name"` // datatypes.JSON
	Count       int64  `json:"count"`
	TotalAmount int64  `json:"total_amount"`
}

func NewJournalArticlePurchaseStatsEntity(journalID uint, journalName []byte, count int64, totalAmount int64) *JournalArticlePurchaseStatsEntity {
	return &JournalArticlePurchaseStatsEntity{JournalID: journalID, JournalName: journalName, Count: count, TotalAmount: totalAmount}
}

func (this *JournalArticlePurchaseStatsEntity) GetJournalName() map[string]string {
	var nameMap map[string]string
	if err := json.Unmarshal(this.JournalName, &nameMap); err != nil {
		nameMap = make(map[string]string)
	}
	return nameMap
}

type JournalArticlePurchaseStatsWithNameEntity struct {
	JournalID   uint              `json:"journal_id"`
	JournalName map[string]string `json:"journal_name"`
	Count       int64             `json:"count"`
	TotalAmount int64             `json:"total_amount"`
}

func NewJournalArticlePurchaseStatsWithNameEntity(journalID uint, journalName map[string]string, count int64, totalAmount int64) *JournalArticlePurchaseStatsWithNameEntity {
	return &JournalArticlePurchaseStatsWithNameEntity{JournalID: journalID, JournalName: journalName, Count: count, TotalAmount: totalAmount}
}
