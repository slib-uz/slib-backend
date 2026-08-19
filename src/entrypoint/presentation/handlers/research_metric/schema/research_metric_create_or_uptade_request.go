package schema

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ResearchMetricCreateOrUpdateRequest struct {
	ProfileUrl string                  `json:"profile_url" validate:"required"`
	HIndex     *uint                   `json:"h_index" validate:"required"`
	Source     enum.ResearchMetricEnum `json:"source" validate:"required"`
}

func (this *ResearchMetricCreateOrUpdateRequest) ToEntity() *entity.ResearchMetricEntity {
	hIndex := uint(0)
	if this.HIndex != nil {
		hIndex = *this.HIndex
	}
	return entity.NewResearchMetricEntity(0, 0, this.ProfileUrl, hIndex, this.Source)
}
