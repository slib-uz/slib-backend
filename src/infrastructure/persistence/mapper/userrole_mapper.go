package mapper

import (
	"encoding/json"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func UserRoleEntityToModel(userRole *entity2.UserRoleEntity) *models.RoleModel {
	return models.NewUserRoleModel(
		userRole.UserID,
		userRole.PublisherID,
		userRole.JournalID,
		userRole.InstitutionID,
		userRole.Role,
		userRole.URL,
	)
}

func UserRoleModelToEntity(userRole *models.RoleModel) *entity2.UserRoleEntity {
	var journalName map[string]string
	var publisherName *string
	var institutionName *string
	var user *entity2.UserEntity
	if userRole.Journal != nil {
		_ = json.Unmarshal(userRole.Journal.Name, &journalName)
	}
	if userRole.Publisher != nil {
		publisherName = userRole.Publisher.Name
	}
	if userRole.Institution != nil {
		institutionName = &userRole.Institution.Name
	}

	if userRole.User != nil {
		user = UserModelToEntity(userRole.User)
	}

	return entity2.NewUserRoleEntity(
		userRole.ID,
		userRole.UserID,
		user,
		userRole.Role,
		userRole.PublisherID,
		userRole.JournalID,
		userRole.InstitutionID,
		publisherName,
		journalName,
		institutionName,
		userRole.URL,
	)
}
