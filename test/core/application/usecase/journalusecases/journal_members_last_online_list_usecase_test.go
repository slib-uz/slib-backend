package journalusecases_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeLastOnlineRoleRepo — UserRoleRepository interfeysi embedded. A'zolar
// ro'yxati testidagi fake'dan alohida, chunki chaqiruvlar hisoblagichi shu
// endpointning o'ziga tegishli bo'lishi kerak.
type fakeLastOnlineRoleRepo struct {
	repository.UserRoleRepository
	calls int
}

func (f *fakeLastOnlineRoleRepo) GetByJournalIDOrderByLastOnline(
	journalID uint,
) ([]*entity.JournalMemberLastOnlineEntity, error) {
	f.calls++
	return []*entity.JournalMemberLastOnlineEntity{}, nil
}

func newLastOnlineUseCase(roleRepo *fakeLastOnlineRoleRepo) *journalusecases.JournalMembersLastOnlineListUsecase {
	return newLastOnlineUseCaseWithJournalRepo(roleRepo, &fakeJournalRepo{})
}

func newLastOnlineUseCaseWithJournalRepo(
	roleRepo *fakeLastOnlineRoleRepo, journalRepo *fakeJournalRepo,
) *journalusecases.JournalMembersLastOnlineListUsecase {
	permission := permissionusecases.NewJournalMemberPermissionUseCase(journalRepo)
	return journalusecases.NewJournalMembersLastOnlineListUsecase(roleRepo, permission)
}

func TestLastOnlineJournalMemberCanListOwnJournal(t *testing.T) {
	roleRepo := &fakeLastOnlineRoleRepo{}
	uc := newLastOnlineUseCase(roleRepo)

	if _, err := uc.Execute(chiefEditorOf(7), 7); err != nil {
		t.Fatalf("o'z jurnali a'zolarini ko'ra olmadi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("ruxsat berildi, lekin repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

// Bayroq bilan berilgan super-admin (users.is_admin, RoleAdmin rolisiz) —
// a'zolar ro'yxatidagi bilan bir xil qoida.
func TestLastOnlineFlagOnlyAdminCanListAnyJournal(t *testing.T) {
	roleRepo := &fakeLastOnlineRoleRepo{}
	uc := newLastOnlineUseCase(roleRepo)

	admin := &entity.UserBasicEntity{ID: 1, IsAdmin: true, Roles: nil}

	if _, err := uc.Execute(admin, 99); err != nil {
		t.Fatalf("bayroq bilan berilgan super-admin rad etildi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("admin uchun repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

// Asosiy IDOR testi: /members yopilgan bo'lsa ham, /members/last-online orqali
// begona jurnalning tahririyati (F.I.SH., science ID, telefon raqami) ochilmasligi kerak.
func TestLastOnlineForeignJournalIsDenied(t *testing.T) {
	roleRepo := &fakeLastOnlineRoleRepo{}
	uc := newLastOnlineUseCase(roleRepo)

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
// a'zolar repozitoriysiga umuman borilmasligi kerak. Xato aynan o'zi
// qaytarilishi tekshiriladi — `allowed, _ :=` deb soddalashtirilsa test yiqiladi.
func TestLastOnlinePermissionErrorDeniesAndSkipsRepository(t *testing.T) {
	dbErr := errors.New("db down")
	roleRepo := &fakeLastOnlineRoleRepo{}
	journalRepo := &fakeJournalRepo{publisherID: 300, err: dbErr}
	uc := newLastOnlineUseCaseWithJournalRepo(roleRepo, journalRepo)

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
func TestLastOnlinePublisherAdminOfOwningPublisherIsGranted(t *testing.T) {
	roleRepo := &fakeLastOnlineRoleRepo{}
	journalRepo := &fakeJournalRepo{publisherID: 300}
	uc := newLastOnlineUseCaseWithJournalRepo(roleRepo, journalRepo)

	if _, err := uc.Execute(publisherAdminOf(300), 7); err != nil {
		t.Fatalf("o'z nashriyoti jurnaliga publisher admin kira olmadi: %v", err)
	}
	if roleRepo.calls != 1 {
		t.Fatalf("ruxsat berildi, lekin repozitoriy 1 marta chaqirilmadi: %d", roleRepo.calls)
	}
}

func TestLastOnlineNilUserIsRejected(t *testing.T) {
	roleRepo := &fakeLastOnlineRoleRepo{}
	uc := newLastOnlineUseCase(roleRepo)

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
