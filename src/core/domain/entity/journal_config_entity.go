package entity

import "time"

type JournalConfigEntity struct {
	ID uint `json:"id"`

	JournalID uint           `json:"journal_id"`
	Journal   *JournalEntity `json:"journal"`

	CreatorID uint        `json:"creator_id"`
	Creator   *UserEntity `json:"creator"`

	WebsiteURL string         `json:"website_url"`
	Conf       map[string]any `json:"conf"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func NewJournalConfigEntity(
	ID uint,
	journalID uint,
	creatorID uint,
	creator *UserEntity,
	journal *JournalEntity,
	websiteURL string,
	conf map[string]any,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *JournalConfigEntity {
	return &JournalConfigEntity{
		ID:         ID,
		JournalID:  journalID,
		CreatorID:  creatorID,
		Creator:    creator,
		Journal:    journal,
		WebsiteURL: websiteURL,
		Conf:       conf,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
