package schema

import (
	"slib.uz/src/core/domain/entity/enum"
)

type UserUpdateRequest struct {
	Photo          *string                  `json:"photo,omitempty"`
	Email          *string                  `json:"email,omitempty"`
	AcademicDegree *enum.AcademicDegreeCode `json:"academic_degree,omitempty"`
	AcademicTitle  *string                  `json:"academic_title,omitempty"`
	ORCIDID        *string                  `json:"orcid_id,omitempty"`
}

type UserUpdateResponse struct {
	Photo          *string                  `json:"photo,omitempty"`
	Email          *string                  `json:"email,omitempty"`
	AcademicDegree *enum.AcademicDegreeCode `json:"academic_degree,omitempty"`
	AcademicTitle  *string                  `json:"academic_title,omitempty"`
	ORCIDID        *string                  `json:"orcid_id,omitempty"`
}

func NewUserUpdateResponse(photo, email *string, academicDegree *enum.AcademicDegreeCode, academicTitle *string, orcidID *string) *UserUpdateResponse {
	return &UserUpdateResponse{
		Photo:          photo,
		Email:          email,
		AcademicDegree: academicDegree,
		AcademicTitle:  academicTitle,
		ORCIDID:        orcidID,
	}
}
