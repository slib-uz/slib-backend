package schema

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type RoleCreateRequest struct {
	UserID        uint          `json:"user_id" validate:"required"`
	PublisherID   *uint         `json:"publisher_id"`
	JournalID     *uint         `json:"journal_id"`
	InstitutionID *uint         `json:"institution_id"`
	Url           *string       `json:"url"`
	Role          enum.UserRole `json:"role"`
}

func (this *RoleCreateRequest) ToEntity() *entity.UserRoleEntity {
	return entity.NewUserRoleInputEntity(
		this.UserID,
		this.PublisherID,
		this.JournalID,
		this.InstitutionID,
		this.Role,
		this.Url,
	)
}
