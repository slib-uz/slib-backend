package schema_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

// SandboxLoginRequest tekshiruvdan o'tmasa, yaroqsiz raqam ham sandbox
// login oqimiga kirib ketardi. Bu test faqat struct'da phone_uz tagi
// borligini tekshiradi — tag olib tashlansa, buzilishi kerak. Xabar
// matnining aynan qanday ekanligi request_validator_test.go zimmasida,
// shuning uchun bu yerda takrorlanmaydi.
func TestSandboxLoginRequestRejectsMalformedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	err := v.Validate(&schema.SandboxLoginRequest{PhoneNumber: "9989012345678", Otp: "123456"})

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
}

func TestSandboxLoginRequestAcceptsWellFormedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	if err := v.Validate(&schema.SandboxLoginRequest{PhoneNumber: "998901234567", Otp: "123456"}); err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
}
