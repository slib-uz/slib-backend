package response

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type UzSciFormsResponse struct {
	Success bool                `json:"success"`
	Data    []UzSciFormResponse `json:"data"`
}

type UzSciFormResponse struct {
	ID             uint      `json:"id"`
	RatingPeriodID uint      `json:"rating_period_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	Required       bool      `json:"required"`
	OrderNo        int       `json:"order_no"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (this *UzSciFormResponse) ToEntity() *entity.UzSciFormEntity {
	return entity.NewUzSciFormEntity(
		this.ID,
		this.RatingPeriodID,
		this.Name,
		this.Description,
		this.Type,
		this.Required,
		this.OrderNo,
		this.CreatedAt,
		this.UpdatedAt,
	)
}
