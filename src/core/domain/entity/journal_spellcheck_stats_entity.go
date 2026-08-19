package entity

import "encoding/json"

type JournalSpellcheckStatsEntity struct {
	JournalID   uint   `json:"journal_id"`
	JournalName []byte `json:"journal_name"` // datatypes.JSON
	Success     int    `json:"success"`
	Failed      int    `json:"failed"`
}

func NewJournalSpellcheckStatsEntity(journalID uint, journalName []byte, success int, failed int) *JournalSpellcheckStatsEntity {
	return &JournalSpellcheckStatsEntity{JournalID: journalID, JournalName: journalName, Success: success, Failed: failed}
}

func (this *JournalSpellcheckStatsEntity) GetJournalName() map[string]string {
	var nameMap map[string]string
	if err := json.Unmarshal(this.JournalName, &nameMap); err != nil {
		nameMap = make(map[string]string)
	}
	return nameMap
}

type JournalSpellcheckStatsWithNameEntity struct {
	JournalID   uint              `json:"journal_id"`
	JournalName map[string]string `json:"journal_name"`
	Success     int               `json:"success"`
	Failed      int               `json:"failed"`
}

func NewJournalSpellcheckStatsWithNameEntity(journalID uint, journalName map[string]string, success int, failed int) *JournalSpellcheckStatsWithNameEntity {
	return &JournalSpellcheckStatsWithNameEntity{JournalID: journalID, JournalName: journalName, Success: success, Failed: failed}
}
