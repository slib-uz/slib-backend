package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func UserRoleModelToEntity(role *models.RoleModel) *entity.UserRoleEntity {
	var journalName map[string]string
	var publisherName *string
	var institutionName *string
	if role.Journal != nil {
		journalName = FromGormJson[map[string]string](role.Journal.Name)
	}
	if role.Publisher != nil {
		publisherName = role.Publisher.Name
	}
	if role.Institution != nil {
		institutionName = &role.Institution.Name
	}
	return entity.NewUserRoleEntity(
		role.ID,
		role.UserID,
		nil,
		role.Role,
		role.PublisherID,
		role.JournalID,
		role.InstitutionID,
		publisherName,
		journalName,
		institutionName,
		role.URL,
	)
}
