package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type UserUpdateEntity struct {
	Photo          *string                  `json:"photo,omitempty"`
	Email          *string                  `json:"email,omitempty"`
	AcademicDegree *enum.AcademicDegreeCode `json:"academic_degree,omitempty"`
	AcademicTitle  *string                  `json:"academic_title,omitempty"`
	ORCIDID        *string                  `json:"orcid_id,omitempty"`
}

func NewUserUpdateEntity(photo, email *string, academicDegree *enum.AcademicDegreeCode, academicTitle *string, orcidID *string) *UserUpdateEntity {
	return &UserUpdateEntity{
		Photo:          photo,
		Email:          email,
		AcademicDegree: academicDegree,
		AcademicTitle:  academicTitle,
		ORCIDID:        orcidID,
	}
}
