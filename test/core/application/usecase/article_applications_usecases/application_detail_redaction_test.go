package article_applications_usecases_test

import (
	"testing"

	usecase "slib.uz/src/core/application/usecase/article_applications_usecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

// Ariza egasi (owner) — testlarda "o'zi" va "begona"ni ajratish uchun.
const applicationOwnerID uint = 10

type fakeApplicationDetailRepo struct {
	repository.ApplicationRepository
	application *entity.ApplicationEntity
}

func (f *fakeApplicationDetailRepo) GetByIDWithRelations(id uint) (*entity.ApplicationEntity, error) {
	return f.application, nil
}

type fakeApplicationDetailReferenceRepo struct {
	repository.ArticleReferenceRepository
}

func (f *fakeApplicationDetailReferenceRepo) GetListByArticleID(articleID uint) ([]*entity.ReferenceEntity, error) {
	return nil, nil
}

// buildAppWithContacts ariza egasi, affiliation muallif va reviewer uchun
// telefon/email bilan to'ldirilgan ApplicationEntity quradi — RedactApplicationContacts
// ning uchala yo'lini ham sinash uchun.
func buildAppWithContacts() *entity.ApplicationEntity {
	return &entity.ApplicationEntity{
		ID:        1,
		ArticleID: 100,
		UserID:    applicationOwnerID,
		User: &entity.UserEntity{
			ID:          applicationOwnerID,
			PhoneNumber: "+998900000010",
			Email:       "owner@x.uz",
		},
		Article: &entity.ArticleEntity{
			ID: 100,
			ArticleAuthorAffiliation: []*entity.ArticleAuthorAffiliationEntity{
				{
					ID: 1,
					Author: &entity.AuthorEntity{
						ID: 20,
						User: &entity.UserEntity{
							ID:          20,
							PhoneNumber: "+998900000020",
							Email:       "aff@x.uz",
						},
					},
				},
			},
		},
		ReviewStages: []*entity.ReviewStageEntity{
			{
				ID: 1,
				Reviewer: &entity.UserEntity{
					ID:          30,
					PhoneNumber: "+998900000030",
					Email:       "reviewer@x.uz",
				},
			},
		},
	}
}

func newApplicationDetailUsecase(app *entity.ApplicationEntity) *usecase.ApplicationDetailUsecase {
	return usecase.NewApplicationDetailUsecase(
		&fakeApplicationDetailRepo{application: app},
		&fakeApplicationDetailReferenceRepo{},
	)
}

// Asosiy talab (CWE-200): /article-application/detail/{id} begona so'rovchiga
// ariza egasi, affiliation muallif va reviewer'ning telefon/emailini ko'rsatmasligi
// kerak.
func TestApplicationDetailRedactsStrangerContacts(t *testing.T) {
	uc := newApplicationDetailUsecase(buildAppWithContacts())

	resp, err := uc.Execute(&entity.UserBasicEntity{ID: 999}, 1)
	if err != nil {
		t.Fatalf("ariza detali o'qilmadi: %v", err)
	}

	if resp.User.PhoneNumber != "" || resp.User.Email != "" {
		t.Errorf("begonaga ega telefon/email ochildi: %q / %q", resp.User.PhoneNumber, resp.User.Email)
	}
	aff := resp.Article.CoAuthorsWithAffiliation[0].Author.User
	if aff.PhoneNumber != "" || aff.Email != "" {
		t.Errorf("begonaga affiliation muallif telefon/email ochildi: %q / %q", aff.PhoneNumber, aff.Email)
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber != "" {
		t.Errorf("begonaga reviewer telefoni ochildi: %q", resp.ReviewStages[0].Reviewer.PhoneNumber)
	}
}

func TestApplicationDetailKeepsOwnerContacts(t *testing.T) {
	uc := newApplicationDetailUsecase(buildAppWithContacts())

	resp, err := uc.Execute(&entity.UserBasicEntity{ID: applicationOwnerID}, 1)
	if err != nil {
		t.Fatalf("ariza detali o'qilmadi: %v", err)
	}

	if resp.User.PhoneNumber == "" || resp.User.Email == "" {
		t.Error("egaga o'z telefoni/emaili ko'rinishi kerak edi")
	}
}

func TestApplicationDetailKeepsAdminContacts(t *testing.T) {
	uc := newApplicationDetailUsecase(buildAppWithContacts())

	admin := &entity.UserBasicEntity{ID: 999, IsAdmin: true}
	resp, err := uc.Execute(admin, 1)
	if err != nil {
		t.Fatalf("ariza detali o'qilmadi: %v", err)
	}

	if resp.User.PhoneNumber == "" || resp.User.Email == "" {
		t.Error("admin uchun telefon/email ko'rinishi kerak edi")
	}
	aff := resp.Article.CoAuthorsWithAffiliation[0].Author.User
	if aff.PhoneNumber == "" || aff.Email == "" {
		t.Error("admin uchun affiliation muallif telefon/email ko'rinishi kerak edi")
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber == "" {
		t.Error("admin uchun reviewer telefoni ko'rinishi kerak edi")
	}
}

// nil so'rovchi (masalan noto'g'ri kontekst) xavfsiz tomonga: hech kim ega bo'lmaydi,
// hammasi bo'shaladi, panika bo'lmaydi.
func TestApplicationDetailNilRequesterBlanksAll(t *testing.T) {
	uc := newApplicationDetailUsecase(buildAppWithContacts())

	resp, err := uc.Execute(nil, 1)
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if resp.User.PhoneNumber != "" || resp.User.Email != "" {
		t.Error("nil so'rovchiga ega telefon/email ochildi")
	}
}
