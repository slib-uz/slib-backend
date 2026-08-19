package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
)

func UserRoleEntityListToWithBasicUserEntityList(userRoles []*entity2.UserRoleEntity) []*entity2.UserRoleWithBasicUserEntity {
	userRoleWithBasicUserEntities := make([]*entity2.UserRoleWithBasicUserEntity, len(userRoles))
	for i, userRole := range userRoles {
		userRoleWithBasicUserEntities[i] = UserRoleEntityToWithBasicUserEntity(userRole)
	}
	return userRoleWithBasicUserEntities
}

func UserRoleEntityToWithBasicUserEntity(userRole *entity2.UserRoleEntity) *entity2.UserRoleWithBasicUserEntity {

	var user *entity2.UserBasicEntity
	if userRole.User != nil {
		user = UserEntityToBasic(userRole.User)
	}

	return entity2.NewUserRoleWithBasicUserEntity(
		userRole.ID,
		userRole.UserID,
		user,
		userRole.Role,
		userRole.PublisherID,
		userRole.JournalID,
		userRole.InstitutionID,
		userRole.PublisherName,
		userRole.JournalName,
		userRole.InstitutionName,
		userRole.URL,
	)
}
