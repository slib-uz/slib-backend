package entity

import "time"

type JournalRatingEntity struct {
	ID        uint                `json:"id"`
	UserID    uint                `json:"user_id"`
	User      *JournalRaterEntity `json:"user"`
	JournalID uint                `json:"journal_id"`
	Stars     uint                `json:"stars"`
	Review    string              `json:"review"`
	CreatedAt time.Time           `json:"created_at"`
}

func NewJournalRatingEntity(ID uint, userID uint, user *JournalRaterEntity, journalID uint, stars uint, review string, createdAt time.Time) *JournalRatingEntity {
	return &JournalRatingEntity{ID: ID, UserID: userID, User: user, JournalID: journalID, Stars: stars, Review: review, CreatedAt: createdAt}
}
