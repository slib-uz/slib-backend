package entity

import (
	"time"

	enum2 "slib.uz/src/core/domain/entity/enum"
)

type UserEntity struct {
	ID uint `json:"id"`

	ScienceID      string                   `json:"science_id"`
	Pin            *string                  `json:"-"`
	FirstName      string                   `json:"first_name"`
	LastName       string                   `json:"last_name"`
	MiddleName     string                   `json:"middle_name"`
	FullName       string                   `json:"full_name"`
	PhoneNumber    string                   `json:"phone_number"`
	IsAdmin        bool                     `json:"is_admin"`
	Gender         uint                     `json:"gender"`
	BirthDate      *time.Time               `json:"-"`
	Roles          []*UserRoleEntity        `json:"roles"`
	Photo          string                   `json:"photo"`
	Email          string                   `json:"email"`
	City           enum2.CityCode           `json:"city" swaggertype:"string"`
	AcademicDegree enum2.AcademicDegreeCode `json:"academic_degree" swaggertype:"string"`
	AcademicTitle  *string                  `json:"academic_title"`
	ORCIDID        *string                  `json:"orcid_id"`
	ArticleCount   int64                    `json:"article_count"`
}

func NewUserEntity(ID uint, scienceID string, pin *string, firstName, lastName, middleName, fullName, phoneNumber string, isAdmin bool, gender uint, birthDate *time.Time, roles []*UserRoleEntity, photo string, email string, city enum2.CityCode, academicDegree enum2.AcademicDegreeCode, academicTitle *string, orcidID *string, articleCount int64) *UserEntity {
	return &UserEntity{
		ID:             ID,
		ScienceID:      scienceID,
		Pin:            pin,
		FirstName:      firstName,
		LastName:       lastName,
		MiddleName:     middleName,
		FullName:       fullName,
		PhoneNumber:    phoneNumber,
		IsAdmin:        isAdmin,
		Gender:         gender,
		BirthDate:      birthDate,
		Roles:          roles,
		Photo:          photo,
		Email:          email,
		City:           city,
		AcademicDegree: academicDegree,
		AcademicTitle:  academicTitle,
		ORCIDID:        orcidID,
		ArticleCount:   articleCount,
	}
}

type UserDTOForAuthorEntity struct {
	ID             uint    `json:"id"`
	ScienceID      string  `json:"science_id"`
	FullName       string  `json:"full_name"`
	AcademicTitle  *string `json:"academic_title"`
	AcademicDegree *string `json:"academic_degree"`
}

func NewUserDTOForAuthorEntity(ID uint, scienceId, fullName string, academicTitle *string, academicDegree *string) *UserDTOForAuthorEntity {
	return &UserDTOForAuthorEntity{ID: ID, ScienceID: scienceId, FullName: fullName, AcademicTitle: academicTitle, AcademicDegree: academicDegree}
}
