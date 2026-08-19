package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

// IsAdmin foydalanuvchining super-admin ekanini aniqlaydi.
//
// Loyihada admin ikki xil ifodalanadi: UserEntity dagi IsAdmin bayrog'i va
// RoleAdmin roli. Ikkalasi ham qabul qilinadi, shunda tekshiruv qaysi
// mexanizm ishlatilganidan qat'i nazar bir xil javob beradi.
func IsAdmin(user *entity.UserBasicEntity) bool {
	if user == nil {
		return false
	}

	if user.IsAdmin {
		return true
	}

	for _, role := range user.Roles {
		if role != nil && role.Role == enum.RoleAdmin {
			return true
		}
	}

	return false
}
