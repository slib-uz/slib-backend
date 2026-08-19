package entity

import "time"

type JournalDoiSettingEntity struct {
	ID          uint
	JournalID   uint
	JournalName string
	Username    string
	Password    string
	DOIPrefix   string // e.g. "10.66701"
	DOISuffix   string // e.g. "ADPIIX"
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewJournalDoiSettingEntity(id uint, journalID uint, journalName string, username string, password string, doiPrefix string, doiSuffix string, createdAt time.Time, updatedAt time.Time) *JournalDoiSettingEntity {
	return &JournalDoiSettingEntity{ID: id, JournalID: journalID, JournalName: journalName, Username: username, Password: password, DOIPrefix: doiPrefix, DOISuffix: doiSuffix, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

// GenerateDOI generates DOI for an article: {prefix}/{suffix}/art{articleID}
func (this *JournalDoiSettingEntity) GenerateDOI(articleID uint) string {
	return this.DOIPrefix + "/" + this.DOISuffix + "/art" + uintToStr(articleID)
}

func uintToStr(v uint) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
