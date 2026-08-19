package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func UserEntityToBasic(user *entity.UserEntity) *entity.UserBasicEntity {
	var roles []*entity.UserRoleEntity
	for _, item := range user.Roles {
		roles = append(roles, item)
	}
	return entity.NewUserBasicEntity(
		user.ID,
		user.ScienceID,
		//user.Pin,
		user.FullName,
		user.PhoneNumber,
		user.IsAdmin,
		roles,
	)
}

func UserEntityToShared(user *entity.UserEntity) *entity.UserSharedEntity {
	var roles []*entity.UserRoleEntity
	for _, item := range user.Roles {
		roles = append(roles, item)
	}
	return entity.NewUserSharedEntity(
		user.ID,
		user.ScienceID,
		user.FullName,
	)
}
