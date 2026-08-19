package entity

import "time"

type JournalEditorialEntity struct {
	ID        uint      `json:"id"`
	JournalID uint      `json:"journal_id"`
	FullName  string    `json:"full_name"`
	RoleCode  int       `json:"role_code"`
	RoleTitle string    `json:"role_title"`
	Photo     string    `json:"photo"`
	ScienceID string    `json:"science_id"`
	Workplace string    `json:"workplace"`
	Position  string    `json:"position"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"-"`
}
