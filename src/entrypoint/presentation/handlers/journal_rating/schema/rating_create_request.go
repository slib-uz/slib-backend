package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type JournalRatingCreateRequest struct {
	JournalID uint   `json:"journal_id"`
	Stars     uint   `json:"stars"`
	Review    string `json:"review"`
}

func (this *JournalRatingCreateRequest) ToEntity() *entity.JournalRatingEntity {
	return entity.NewJournalRatingEntity(0, 0, nil, this.JournalID, this.Stars, this.Review, time.Now())
}
