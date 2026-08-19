package models

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type UserModel struct {
	gorm.Model

	ScienceID string  `gorm:"size:64;uniqueIndex;not null;"` // public
	Pin       *string `gorm:"size:64;uniqueIndex;"`

	FirstName      string `gorm:"size:255;"`
	LastName       string `gorm:"size:255;"`
	MiddleName     string `gorm:"size:255;"`
	FullName       string `gorm:"size:255;"` // public
	Gender         uint
	BirthDate      *time.Time              `gorm:"type:date"`
	PhoneNumber    string                  `gorm:"size:64;;not null;uniqueIndex;"`
	IsAdmin        bool                    `gorm:"default:false"`
	Roles          []*RoleModel            `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL;"`
	Photo          string                  `gorm:"size:255;"`
	Email          string                  `gorm:"size:255;"`
	City           enum.CityCode           `gorm:"size:5;not null;default:''"`
	AcademicDegree enum.AcademicDegreeCode `gorm:"size:5;not null;default:''"`
	AcademicTitle  *string                 `gorm:"size:255;"`
	ORCIDID        *string                 `gorm:"size:50;uniqueIndex;"`
}

func NewUserModel(scienceID string, pin *string, firstName string, lastName string, middleName string, fullName string, gender uint, birthDate *time.Time, phoneNumber string, photo string, email string, city enum.CityCode, academicDegree enum.AcademicDegreeCode, academicTitle *string, orcidID *string) *UserModel {
	return &UserModel{ScienceID: scienceID, Pin: pin, FirstName: firstName, LastName: lastName, MiddleName: middleName, FullName: fullName, Gender: gender, BirthDate: birthDate, PhoneNumber: phoneNumber, Photo: photo, Email: email, City: city, AcademicDegree: academicDegree, AcademicTitle: academicTitle, ORCIDID: orcidID}
}

func (UserModel) TableName() string {
	return "users"
}
