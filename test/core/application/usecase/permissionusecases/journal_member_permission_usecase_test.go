package permissionusecases_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeJournalRepo — JournalRepository interfeysi embedded, faqat kerakli metod yoziladi.
type fakeJournalRepo struct {
	repository.JournalRepository
	publisherID uint
	err         error
	calls       int
}

func (f *fakeJournalRepo) GetPublisherIdByJournalId(journalID uint) (uint, error) {
	f.calls++
	return f.publisherID, f.err
}

func journalRole(role enum.UserRole, journalID uint) *entity.UserRoleEntity {
	id := journalID
	return &entity.UserRoleEntity{Role: role, JournalID: &id}
}

func publisherRole(publisherID uint) *entity.UserRoleEntity {
	id := publisherID
	return &entity.UserRoleEntity{Role: enum.RolePublisherAdmin, PublisherID: &id}
}

func TestChiefEditorAllowedOnOwnJournal(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleChiefEditor, 7)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("bosh muharrir o'z jurnaliga kira olmadi")
	}
	if repo.calls != 0 {
		t.Fatalf("keraksiz DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestChiefEditorDeniedOnOtherJournal(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleChiefEditor, 7)}, 99)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("bosh muharrir begona jurnalga kirdi")
	}
}

// Zaiflikning o'zi: hozir istalgan nashriyot admini istalgan jurnalga kiradi.
func TestPublisherAdminDeniedOnOtherPublishersJournal(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 500} // jurnal 500-nashriyotga tegishli
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("begona nashriyot admini jurnalga kirdi — IDOR ochiq")
	}
}

func TestPublisherAdminAllowedOnOwnPublishersJournal(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 300}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("o'z nashriyoti jurnaliga publisher admin kira olmadi")
	}
}

func TestAdminAlwaysAllowedWithoutLookup(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{{Role: enum.RoleAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("admin rad etildi")
	}
	if repo.calls != 0 {
		t.Fatalf("admin uchun keraksiz DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestNoLookupWhenUserHasNoPublisherAdminRole(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	_, err := uc.Execute([]*entity.UserRoleEntity{journalRole(enum.RoleSecretary, 99)}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if repo.calls != 0 {
		t.Fatalf("publisher-admin roli yo'q, lekin DB so'rovi ketdi: %d", repo.calls)
	}
}

func TestRepositoryErrorIsPropagatedAndDenies(t *testing.T) {
	repo := &fakeJournalRepo{err: errors.New("db down")}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{publisherRole(300)}, 7)
	if err == nil {
		t.Fatal("repozitoriy xatosi yutib yuborildi")
	}
	if allowed {
		t.Fatal("xato holatida ruxsat berildi — fail-closed buzildi")
	}
}

func TestNilRolesAreSkipped(t *testing.T) {
	repo := &fakeJournalRepo{}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{nil, {Role: enum.RoleAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if !allowed {
		t.Fatal("nil rol butun ro'yxatni buzdi")
	}
}

func TestPublisherRoleWithNilPublisherIDDoesNotPanic(t *testing.T) {
	repo := &fakeJournalRepo{publisherID: 300}
	uc := permissionusecases.NewJournalMemberPermissionUseCase(repo)

	allowed, err := uc.Execute([]*entity.UserRoleEntity{{Role: enum.RolePublisherAdmin}}, 7)
	if err != nil {
		t.Fatalf("xato qaytdi: %v", err)
	}
	if allowed {
		t.Fatal("PublisherID nil bo'lgan rol ruxsat berdi")
	}
}
