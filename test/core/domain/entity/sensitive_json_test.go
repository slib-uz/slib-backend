package entity_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"slib.uz/src/core/domain/entity"
)

// PINFL va tug'ilgan sana hech qanday API javobida chiqmasligi kerak (CWE-200).
func TestUserEntityDoesNotExposePinOrBirthDate(t *testing.T) {
	pin := "12345678901234"
	birth := time.Now()
	u := &entity.UserEntity{
		ID:          1,
		Pin:         &pin,
		BirthDate:   &birth,
		PhoneNumber: "+998901234567",
		Email:       "a@b.uz",
	}

	raw, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal xatosi: %v", err)
	}
	out := string(raw)

	if strings.Contains(out, "\"pin\"") {
		t.Error("javobda \"pin\" kaliti bor, bo'lmasligi kerak")
	}
	if strings.Contains(out, "\"birth_date\"") {
		t.Error("javobda \"birth_date\" kaliti bor, bo'lmasligi kerak")
	}
	if strings.Contains(out, pin) {
		t.Error("javobda PINFL qiymati bor")
	}
	// Telefon bu darajada qoladi — u so'rovchiga qarab use case'da redaksiya qilinadi.
	if !strings.Contains(out, "phone_number") {
		t.Error("phone_number kaliti bo'lishi kerak edi (bu darajada redaksiya yo'q)")
	}
}

// Reviewer entity hech qanday API javobida telefon raqamini oshkor qilmasligi kerak (CWE-200).
func TestReviewerEntityDoesNotExposePhone(t *testing.T) {
	r := &entity.ReviewerEntity{ID: 1, FullName: "A B", PhoneNumber: "+998901234567"}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal xatosi: %v", err)
	}
	if strings.Contains(string(raw), "phone_number") {
		t.Error("reviewer javobida phone_number bor, bo'lmasligi kerak")
	}
	if strings.Contains(string(raw), "+998901234567") {
		t.Error("reviewer javobida telefon qiymati bor")
	}
}
