package entity

import "encoding/json"

type JournalAntiPlagStatsEntity struct {
	JournalID   uint   `json:"journal_id"`
	JournalName []byte `json:"journal_name"`
	Success     int    `json:"success"`
	Failed      int    `json:"failed"`
}

func NewJournalAntiPlagStatsEntity(journalID uint, journalName []byte, success int, failed int) *JournalAntiPlagStatsEntity {
	return &JournalAntiPlagStatsEntity{JournalID: journalID, JournalName: journalName, Success: success, Failed: failed}
}

func (this *JournalAntiPlagStatsEntity) GetJournalName() map[string]string {
	var nameMap map[string]string
	if err := json.Unmarshal(this.JournalName, &nameMap); err != nil {
		nameMap = make(map[string]string)
	}
	return nameMap
}

type JournalAntiPlagStatsWithNameEntity struct {
	JournalID   uint              `json:"journal_id"`
	JournalName map[string]string `json:"journal_name"`
	Success     int               `json:"success"`
	Failed      int               `json:"failed"`
}

func NewJournalAntiPlagStatsWithNameEntity(journalID uint, journalName map[string]string, success int, failed int) *JournalAntiPlagStatsWithNameEntity {
	return &JournalAntiPlagStatsWithNameEntity{JournalID: journalID, JournalName: journalName, Success: success, Failed: failed}
}
