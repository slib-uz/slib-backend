package permissionusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
)

func TestRedactUserContactBlanksForStranger(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 99, false)
	if u.PhoneNumber != "" || u.Email != "" {
		t.Errorf("begona uchun telefon/email bo'sh bo'lishi kerak, %q / %q keldi", u.PhoneNumber, u.Email)
	}
}

func TestRedactUserContactKeepsForOwner(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 5, false)
	if u.PhoneNumber == "" || u.Email == "" {
		t.Error("egaga o'z aloqa ma'lumoti ko'rinishi kerak edi")
	}
}

func TestRedactUserContactKeepsForPrivileged(t *testing.T) {
	u := &entity.UserEntity{ID: 5, PhoneNumber: "+998901112233", Email: "a@b.uz"}
	permissionusecases.RedactUserContact(u, 99, true)
	if u.PhoneNumber == "" || u.Email == "" {
		t.Error("admin uchun aloqa ma'lumoti ko'rinishi kerak edi")
	}
}

func TestRedactUserContactNilSafe(t *testing.T) {
	permissionusecases.RedactUserContact(nil, 1, false) // panic bermasligi kerak
}

func TestRedactBasicContactBlanksForStranger(t *testing.T) {
	u := &entity.UserBasicEntity{ID: 5, PhoneNumber: "+998901112233"}
	permissionusecases.RedactBasicContact(u, 99, false)
	if u.PhoneNumber != "" {
		t.Errorf("begona uchun telefon bo'sh bo'lishi kerak, %q keldi", u.PhoneNumber)
	}
}

func TestRedactBasicContactKeepsForOwner(t *testing.T) {
	u := &entity.UserBasicEntity{ID: 5, PhoneNumber: "+998901112233"}
	permissionusecases.RedactBasicContact(u, 5, false)
	if u.PhoneNumber == "" {
		t.Error("egaga o'z telefoni ko'rinishi kerak edi")
	}
}

func TestRedactBasicContactNilSafe(t *testing.T) {
	permissionusecases.RedactBasicContact(nil, 1, false) // panic bermasligi kerak
}

// RedactApplicationContacts javobning uchala maxfiy yo'lini birdan tozalaydi:
// ega (.User), affiliation muallif (.Article.CoAuthorsWithAffiliation[].Author.User),
// va reviewer (.ReviewStages[].Reviewer). Fake repo kerak emas — response entity
// to'g'ridan-to'g'ri quriladi.
func TestRedactApplicationContactsBlanksStranger(t *testing.T) {
	resp := &entity.ApplicationResponseEntity{
		User: &entity.UserEntity{ID: 10, PhoneNumber: "+998900000010", Email: "owner@x.uz"},
		Article: &entity.ArticleInputEntity{
			CoAuthorsWithAffiliation: []*entity.ArticleAuthorAffiliationEntity{
				{Author: &entity.AuthorEntity{User: &entity.UserEntity{ID: 20, PhoneNumber: "+998900000020", Email: "aff@x.uz"}}},
			},
		},
		ReviewStages: []*entity.ReviewStageResponseEntity{
			{Reviewer: &entity.UserBasicEntity{ID: 30, PhoneNumber: "+998900000030"}},
		},
	}

	permissionusecases.RedactApplicationContacts(resp, 999, false)

	if resp.User.PhoneNumber != "" || resp.User.Email != "" {
		t.Error("ega telefon/email bo'sh bo'lishi kerak")
	}
	if aff := resp.Article.CoAuthorsWithAffiliation[0].Author.User; aff.PhoneNumber != "" || aff.Email != "" {
		t.Error("affiliation muallif telefon/email bo'sh bo'lishi kerak")
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber != "" {
		t.Error("reviewer telefon bo'sh bo'lishi kerak")
	}
}

func TestRedactApplicationContactsKeepsOwnerRecord(t *testing.T) {
	// So'rovchi = ega (ID 10). Eganing o'z yozuvi ko'rinadi, begonalarniki emas.
	resp := &entity.ApplicationResponseEntity{
		User: &entity.UserEntity{ID: 10, PhoneNumber: "+998900000010", Email: "owner@x.uz"},
		ReviewStages: []*entity.ReviewStageResponseEntity{
			{Reviewer: &entity.UserBasicEntity{ID: 30, PhoneNumber: "+998900000030"}},
		},
	}

	permissionusecases.RedactApplicationContacts(resp, 10, false)

	if resp.User.PhoneNumber == "" {
		t.Error("egaga o'z telefoni ko'rinishi kerak edi")
	}
	if resp.ReviewStages[0].Reviewer.PhoneNumber != "" {
		t.Error("begona reviewer telefoni bo'sh bo'lishi kerak")
	}
}

func TestRedactApplicationContactsNilSafe(t *testing.T) {
	permissionusecases.RedactApplicationContacts(nil, 1, false) // panic bermasligi kerak
	// Article nil va ReviewStages nil holatlar ham xavfsiz:
	permissionusecases.RedactApplicationContacts(&entity.ApplicationResponseEntity{}, 1, false)
}
