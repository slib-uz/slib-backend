package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalMemberPermissionUseCase struct {
	journalRepository repository.JournalRepository
}

// @inject
func NewJournalMemberPermissionUseCase(journalRepository repository.JournalRepository) *JournalMemberPermissionUseCase {
	return &JournalMemberPermissionUseCase{journalRepository: journalRepository}
}

// Execute foydalanuvchining jurnal boshqaruviga kirish huquqini tekshiradi.
//
// Ruxsat beriladigan holatlar:
//   - RoleAdmin — har doim;
//   - RoleChiefEditor yoki RoleSecretary — faqat o'sha jurnalda;
//   - RolePublisherAdmin — faqat jurnal egasi bo'lgan nashriyotda.
//
// Nashriyot egaligini aniqlash uchun DB so'rovi faqat foydalanuvchida
// RolePublisherAdmin roli bo'lsa va arzon tekshiruvlar natija bermagan
// bo'lsagina bajariladi.
func (this *JournalMemberPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, journalID uint) (bool, error) {
	hasPublisherAdminRole := false

	for _, role := range userRoles {
		if role == nil {
			continue
		}

		if role.Role == enum.RoleAdmin {
			return true, nil
		}

		isMemberRole := role.Role == enum.RoleChiefEditor || role.Role == enum.RoleSecretary
		if isMemberRole && role.JournalID != nil && *role.JournalID == journalID {
			return true, nil
		}

		if role.Role == enum.RolePublisherAdmin && role.PublisherID != nil {
			hasPublisherAdminRole = true
		}
	}

	if !hasPublisherAdminRole {
		return false, nil
	}

	ownerPublisherID, err := this.journalRepository.GetPublisherIdByJournalId(journalID)
	if err != nil {
		return false, err
	}

	for _, role := range userRoles {
		if role == nil {
			continue
		}
		if role.Role == enum.RolePublisherAdmin && role.PublisherID != nil && *role.PublisherID == ownerPublisherID {
			return true, nil
		}
	}

	return false, nil
}
