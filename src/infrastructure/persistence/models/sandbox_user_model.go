package models

import "gorm.io/gorm"

type SandboxUserModel struct {
	gorm.Model

	FullName    string
	ScienceID   string
	PhoneNumber string
	Otp         string
}

func (SandboxUserModel) TableName() string {
	return "sandbox_users"
}
