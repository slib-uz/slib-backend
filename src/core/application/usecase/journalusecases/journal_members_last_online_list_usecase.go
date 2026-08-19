package journalusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalMembersLastOnlineListUsecase struct {
	repository       repository.UserRoleRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewJournalMembersLastOnlineListUsecase(
	repository repository.UserRoleRepository,
	memberPermission *permissionusecases.JournalMemberPermissionUseCase,
) *JournalMembersLastOnlineListUsecase {
	return &JournalMembersLastOnlineListUsecase{repository: repository, memberPermission: memberPermission}
}

// Execute jurnal tahririyat a'zolarini oxirgi faollik vaqti bilan qaytaradi.
//
// Javob a'zolar ro'yxati bilan bir xil shaxsiy maydonlarni (F.I.SH., science ID,
// telefon raqami) va qo'shimcha ravishda oxirgi kirish vaqtini o'z ichiga oladi,
// shuning uchun ruxsat nazorati JournalMembersListUsecase bilan bir xil bo'lishi
// shart: aks holda bu endpoint himoyalangan ro'yxatni chetlab o'tish yo'liga
// aylanadi.
func (this *JournalMembersLastOnlineListUsecase) Execute(
	user *entity.UserBasicEntity, journalID uint,
) ([]*entity.JournalMemberLastOnlineEntity, error) {

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

	return this.repository.GetByJournalIDOrderByLastOnline(journalID)
}
