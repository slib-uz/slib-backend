package permissionusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

func TestIsAdminReturnsFalseForNilUser(t *testing.T) {
	if permissionusecases.IsAdmin(nil) {
		t.Fatal("nil foydalanuvchi admin deb topildi")
	}
}

func TestIsAdminHonoursFlag(t *testing.T) {
	user := &entity.UserBasicEntity{ID: 1, IsAdmin: true}
	if !permissionusecases.IsAdmin(user) {
		t.Fatal("IsAdmin bayrog'i e'tiborga olinmadi")
	}
}

func TestIsAdminHonoursRole(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleAdmin}},
	}
	if !permissionusecases.IsAdmin(user) {
		t.Fatal("RoleAdmin roli e'tiborga olinmadi")
	}
}

func TestIsAdminRejectsOrdinaryUser(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleSecretary}},
	}
	if permissionusecases.IsAdmin(user) {
		t.Fatal("oddiy foydalanuvchi admin deb topildi")
	}
}

func TestIsAdminSkipsNilRoles(t *testing.T) {
	user := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{nil, {Role: enum.RoleAdmin}},
	}
	if !permissionusecases.IsAdmin(user) {
		t.Fatal("nil rol butun ro'yxatni buzdi")
	}
}
