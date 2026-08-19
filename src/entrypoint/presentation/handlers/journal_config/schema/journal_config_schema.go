package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type JournalConfigSchema struct {
	JournalID  uint           `json:"journal_id"`
	WebsiteURL string         `json:"website_url"`
	Conf       map[string]any `json:"conf"`
}

func (this *JournalConfigSchema) ToEntity(creatorID uint) *entity.JournalConfigEntity {
	return entity.NewJournalConfigEntity(
		0,
		this.JournalID,
		creatorID,
		nil,
		nil,
		this.WebsiteURL,
		this.Conf,
		true,
		time.Now(),
		time.Now(),
	)
}
