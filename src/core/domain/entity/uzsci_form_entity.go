package entity

import "time"

type UzSciFormEntity struct {
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

func NewUzSciFormEntity(
	id uint,
	ratingPeriodID uint,
	name string,
	description string,
	formType string,
	required bool,
	orderNo int,
	createdAt time.Time,
	updatedAt time.Time,
) *UzSciFormEntity {
	return &UzSciFormEntity{
		ID:             id,
		RatingPeriodID: ratingPeriodID,
		Name:           name,
		Description:    description,
		Type:           formType,
		Required:       required,
		OrderNo:        orderNo,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
