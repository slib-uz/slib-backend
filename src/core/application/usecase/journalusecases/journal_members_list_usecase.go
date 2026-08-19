package journalusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalMembersListUsecase struct {
	repository       repository.UserRoleRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalMembersListUsecase(
	repository repository.UserRoleRepository,
	memberPermission *permissionusecases.JournalMemberPermissionUseCase,
) *JournalMembersListUsecase {
	return &JournalMembersListUsecase{repository: repository, memberPermission: memberPermission}
}

// Execute jurnal tahririyat a'zolari ro'yxatini qaytaradi.
//
// Javob a'zolarning shaxsiy maydonlarini o'z ichiga oladi, shuning uchun
// ro'yxat faqat jurnal boshqaruviga ochiq: bosh muharrir, kotib, jurnal
// egasi bo'lgan nashriyot admini va super-admin.
func (this *JournalMembersListUsecase) Execute(
	user *entity.UserBasicEntity, journalID uint,
) ([]*entity.UserRoleWithBasicUserEntity, error) {

	if user == nil {
		return nil, response.UnauthorizedError
	}

	// Super-admin loyihada ikki xil ifodalanadi: users.is_admin bayrog'i va
	// RoleAdmin roli. JournalMemberPermissionUseCase faqat rolni biladi,
	// shuning uchun bayroq bilan berilgan adminni bu yerda tekshiramiz —
	// aks holda u o'zi ilgari ko'ra olgan ro'yxatdan 403 oladi.
	if !permissionusecases.IsAdmin(user) {
		allowed, err := this.memberPermission.Execute(user.Roles, journalID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, response.PermissionDeniedError
		}
	}

	_members, err := this.repository.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}

	members := make([]*entity.UserRoleEntity, len(_members))
	copy(members, _members)

	return mapper.UserRoleEntityListToWithBasicUserEntityList(members), nil
}
