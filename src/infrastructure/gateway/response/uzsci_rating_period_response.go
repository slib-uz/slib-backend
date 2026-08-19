package response

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type UzSciRatingPeriodsResponse struct {
	Success bool                        `json:"success"`
	Data    []UzSciRatingPeriodResponse `json:"data"`
}

type UzSciRatingPeriodResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Year          int       `json:"year"`
	IsActive      bool      `json:"is_active"`
	MemberGroupID uint      `json:"member_group_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (this *UzSciRatingPeriodResponse) ToEntity() *entity.UzSciRatingPeriodEntity {
	return entity.NewUzSciRatingPeriodEntity(
		this.ID,
		this.Name,
		this.Year,
		this.IsActive,
		this.MemberGroupID,
		this.CreatedAt,
	)
}
