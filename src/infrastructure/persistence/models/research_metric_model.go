package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type ResearchMetricModel struct {
	gorm.Model

	UserID uint       `gorm:"not null;uniqueIndex:idx_user_source"`
	User   *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	ProfileUrl string                  `gorm:"size:255"`
	HIndex     uint                    `gorm:"not null;default:0"`
	Source     enum.ResearchMetricEnum `gorm:"size:32;not null;uniqueIndex:idx_user_source"`
}

func (ResearchMetricModel) TableName() string {
	return "research_metrics"
}

func NewResearchMetricModel(userID uint, profileUrl string, hIndex uint, source enum.ResearchMetricEnum) *ResearchMetricModel {
	return &ResearchMetricModel{UserID: userID, ProfileUrl: profileUrl, HIndex: hIndex, Source: source}
}
