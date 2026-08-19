package entity

import "time"

type EditionEntity struct {
	ID        uint           `json:"id"`
	JournalID uint           `json:"journal_id"`
	Journal   *JournalEntity `json:"journal"`

	File        string    `json:"file"`
	CoverImage  string    `json:"cover_image"`
	Name        string    `json:"name"`
	Number      string    `json:"issue"`
	Volume      string    `json:"volume"`
	PublishedAt time.Time `json:"published_at"`
}

func NewEditionEntity(ID uint, journalID uint, journal *JournalEntity, file string, coverImage string, name string, number string, volume string, publishedAt time.Time) *EditionEntity {
	return &EditionEntity{ID: ID, JournalID: journalID, Journal: journal, File: file, CoverImage: coverImage, Name: name, Number: number, Volume: volume, PublishedAt: publishedAt}
}
