package entity

type JournalRatingStatsEntity struct {
	Total   uint `json:"total"`
	Rating1 uint `json:"rating1"`
	Rating2 uint `json:"rating2"`
	Rating3 uint `json:"rating3"`
	Rating4 uint `json:"rating4"`
	Rating5 uint `json:"rating5"`
}

func NewJournalRatingStatsEntity(total uint, rating1 uint, rating2 uint, rating3 uint, rating4 uint, rating5 uint) *JournalRatingStatsEntity {
	return &JournalRatingStatsEntity{Total: total, Rating1: rating1, Rating2: rating2, Rating3: rating3, Rating4: rating4, Rating5: rating5}
}
