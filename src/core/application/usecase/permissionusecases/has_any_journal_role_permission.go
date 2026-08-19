package permissionusecases

import "slib.uz/src/core/domain/entity"

type HasAnyJournalRolePermissionUseCase struct {
}

// @inject
func NewHasAnyJournalRolePermissionUseCase() *HasAnyJournalRolePermissionUseCase {
	return &HasAnyJournalRolePermissionUseCase{}
}

func (this *HasAnyJournalRolePermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, journalID uint) bool {
	for _, role := range userRoles {
		if role.JournalID != nil && *role.JournalID == journalID {
			return true
		}
	}
	return false
}
