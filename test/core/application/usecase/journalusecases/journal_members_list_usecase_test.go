package journalusecases_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRoleRepo — UserRoleRepository interfeysi embedded.
type fakeUserRoleRepo struct {
	repository.UserRoleRepository
	calls int
}

func (f *fakeUserRoleRepo) GetByJournalID(journalID uint) ([]*entity.UserRoleEntity, error) {
	f.calls++
	return []*entity.UserRoleEntity{}, nil
}

// fakeJournalRepo — JournalRepository interfeysi embedded.
type fakeJournalRepo struct {
	repository.JournalRepository
	publisherID uint
	err         error
}

func (f *fakeJournalRepo) GetPublisherIdByJournalId(journalID uint) (uint, error) {
	return f.publisherID, f.err
}

func chiefEditorOf(journalID uint) *entity.UserBasicEntity {
	id := journalID
	return &entity.UserBasicEntity{
		ID:    42,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleChiefEditor, JournalID: &id}},
	}
}

func publisherAdminOf(publisherID uint) *entity.UserBasicEntity {
	id := publisherID
	return &entity.UserBasicEntity{
		ID:    42,
		Roles: []*entity.UserRoleEntity{{Role: enum.RolePublisherAdmin, PublisherID: &id}},
	}
}

func newUseCase(roleRepo *fakeUserRoleRepo) *journalusecases.JournalMembersListUsecase {
	return newUseCaseWithJournalRepo(roleRepo, &fakeJournalRepo{})
}

func newUseCaseWithJournalRepo(
	roleRepo *fakeUserRoleRepo, journalRepo *fakeJournalRepo,
) *journalusecases.JournalMembersListUsecase {
	permission := permissionusecases.NewJournalMemberPermissionUseCase(journalRepo)
	return journalusecases.NewJournalMembersListUsecase(roleRepo, permission)
}

func TestJournalMemberCanListOwnJournalMembers(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	if _, err := uc.Execute(chiefEditorOf(7), 7); err != nil {
		t.Fatalf("o'z jurnali a'zolarini ko'ra olmadi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("ruxsat berildi, lekin repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

// Super-admin loyihada ikki xil ifodalanadi. AdminPermission middleware faqat
// users.is_admin bayrog'ini qabul qiladi, ya'ni faqat bayrog'i bor adminlar
// prod'da mavjud — ular RoleAdmin rolisiz ham ruxsat olishi shart.
func TestFlagOnlyAdminCanListAnyJournalMembers(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	admin := &entity.UserBasicEntity{ID: 1, IsAdmin: true, Roles: nil}

	if _, err := uc.Execute(admin, 99); err != nil {
		t.Fatalf("bayroq bilan berilgan super-admin rad etildi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("admin uchun repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

func TestForeignJournalMembersAreDenied(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	_, err := uc.Execute(chiefEditorOf(7), 99)
	if err == nil {
		t.Fatal("begona jurnal a'zolari ochildi — IDOR")
	}
	if err != response.PermissionDeniedError {
		t.Fatalf("403 o'rniga boshqa xato qaytdi: %v", err)
	}
	if roleRepo.calls != 0 {
		t.Fatalf("rad etilgan so'rov uchun DB'ga borildi: %d", roleRepo.calls)
	}
}

// Fail-closed: ruxsat tekshiruvi DB xatosiga uchraganda so'rov rad etilishi va
// a'zolar repozitoriysiga umuman borilmasligi kerak.
//
// Xato aynan o'zi qaytarilishi tekshiriladi, PermissionDeniedError emas:
// `allowed, _ :=` deb "soddalashtirilsa" test yiqilishi kerak. Aks holda DB
// uzilishi foydalanuvchiga "sizda ruxsat yo'q" deb yolg'on 403 bo'lib ko'rinadi
// va bu xatolik JournalMemberPermissionUseCase ning fail-closed odatiga
// jimgina bog'lanib qoladi.
func TestPermissionErrorDeniesAndSkipsRepository(t *testing.T) {
	dbErr := errors.New("db down")
	roleRepo := &fakeUserRoleRepo{}
	journalRepo := &fakeJournalRepo{publisherID: 300, err: dbErr}
	uc := newUseCaseWithJournalRepo(roleRepo, journalRepo)

	_, err := uc.Execute(publisherAdminOf(300), 7)
	if err == nil {
		t.Fatal("ruxsat tekshiruvi xatosi yutib yuborildi — xato holatida ruxsat berildi")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("asl DB xatosi o'rniga boshqa xato qaytdi (yutib yuborildi?): %v", err)
	}
	if roleRepo.calls != 0 {
		t.Fatalf("xato holatida a'zolar repozitoriysiga borildi: %d", roleRepo.calls)
	}
}

// Xatosiz holatda o'z nashriyoti jurnaliga publisher admin kira oladi —
// yuqoridagi test faqat xato tufayli rad etilganini tasdiqlash uchun.
func TestPublisherAdminOfOwningPublisherIsGranted(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	journalRepo := &fakeJournalRepo{publisherID: 300}
	uc := newUseCaseWithJournalRepo(roleRepo, journalRepo)

	if _, err := uc.Execute(publisherAdminOf(300), 7); err != nil {
		t.Fatalf("o'z nashriyoti jurnaliga publisher admin kira olmadi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("ruxsat berildi, lekin repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

func TestNilUserIsRejected(t *testing.T) {
	roleRepo := &fakeUserRoleRepo{}
	uc := newUseCase(roleRepo)

	_, err := uc.Execute(nil, 7)
	if err == nil {
		t.Fatal("nil foydalanuvchi qabul qilindi")
	}
	if err != response.UnauthorizedError {
		t.Fatalf("401 o'rniga boshqa xato qaytdi: %v", err)
	}
	if roleRepo.calls != 0 {
		t.Fatalf("nil foydalanuvchi uchun DB'ga borildi: %d", roleRepo.calls)
	}
}
