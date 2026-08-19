package entity

import "time"

type JournalNewsEntity struct {
	ID        uint              `json:"id"`
	JournalID uint              `json:"journal_id"`
	Title     map[string]string `json:"title"`
	Body      map[string]string `json:"body"`
	Image     *string           `json:"image"`
	CreatedAt time.Time         `json:"created_at"`
}
