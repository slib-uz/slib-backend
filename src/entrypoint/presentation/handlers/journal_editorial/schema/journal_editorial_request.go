package schema

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalEditorialCreateOrUpdateRequest struct {
	FullName  string `json:"full_name"`
	RoleCode  int    `json:"role_code"`
	Photo     string `json:"photo"`
	ScienceID string `json:"science_id"`
	Workplace string `json:"workplace"`
	Position  string `json:"position"`
	Order     int    `json:"order"`
}

func (this *JournalEditorialCreateOrUpdateRequest) ToEntity() (*entity.JournalEditorialEntity, error) {
	role, ok := enum.GetJournalEditorialRoleByCode(this.RoleCode)
	if !ok {
		return nil, response.InvalidCodeError
	}

	return &entity.JournalEditorialEntity{
		FullName:  this.FullName,
		RoleCode:  this.RoleCode,
		RoleTitle: role.Label,
		Photo:     this.Photo,
		ScienceID: this.ScienceID,
		Workplace: this.Workplace,
		Position:  this.Position,
		Order:     this.Order,
	}, nil
}
