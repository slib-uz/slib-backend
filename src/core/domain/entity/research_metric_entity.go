package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type ResearchMetricEntity struct {
	ID         uint                    `json:"id"`
	UserID     uint                    `json:"user_id"`
	ProfileUrl string                  `json:"profile_url"`
	HIndex     uint                    `json:"h_index"`
	Source     enum.ResearchMetricEnum `json:"source"`
}

func NewResearchMetricEntity(id, userID uint, profileUrl string, hIndex uint, source enum.ResearchMetricEnum) *ResearchMetricEntity {
	return &ResearchMetricEntity{ID: id, UserID: userID, ProfileUrl: profileUrl, HIndex: hIndex, Source: source}
}
