package entity

import "time"

type UzSciRatingPeriodEntity struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Year          int       `json:"year"`
	IsActive      bool      `json:"is_active"`
	MemberGroupID uint      `json:"member_group_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewUzSciRatingPeriodEntity(
	id uint,
	name string,
	year int,
	isActive bool,
	memberGroupID uint,
	createdAt time.Time,
) *UzSciRatingPeriodEntity {
	return &UzSciRatingPeriodEntity{
		ID:            id,
		Name:          name,
		Year:          year,
		IsActive:      isActive,
		MemberGroupID: memberGroupID,
		CreatedAt:     createdAt,
	}
}
