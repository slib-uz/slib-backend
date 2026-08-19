package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type JournalMemberLastOnlineEntity struct {
	RoleID       uint          `json:"role_id"`
	UserID       uint          `json:"user_id"`
	FullName     string        `json:"full_name"`
	ScienceID    string        `json:"science_id"`
	PhoneNumber  string        `json:"phone_number"`
	Role         enum.UserRole `json:"role"`
	LastOnlineAt *time.Time    `json:"last_online_at"`
}

func NewJournalMemberLastOnlineEntity(roleID, userID uint, fullName, scienceID, phoneNumber string, role enum.UserRole, lastOnlineAt *time.Time) *JournalMemberLastOnlineEntity {
	return &JournalMemberLastOnlineEntity{
		RoleID:       roleID,
		UserID:       userID,
		FullName:     fullName,
		ScienceID:    scienceID,
		PhoneNumber:  phoneNumber,
		Role:         role,
		LastOnlineAt: lastOnlineAt,
	}
}
