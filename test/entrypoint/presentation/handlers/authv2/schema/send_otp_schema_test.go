package schema_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/handlers/authv2/schema"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

// SendOtpRequest tekshiruvdan o'tmasa, yaroqsiz raqam throttle limitini
// sarflab, DB'ga OTP yozib, SMS provayderiga so'rov jo'natardi — o'sha
// provayder 503 qaytarib, mijozga 500 bo'lib yetardi.
//
// Bu test faqat struct'da phone_uz tagi borligini tekshiradi — tag olib
// tashlansa, buzilishi kerak. Xabar matnining aynan qanday ekanligi
// request_validator_test.go zimmasida, shuning uchun bu yerda takrorlanmaydi.
func TestSendOtpRequestRejectsMalformedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	err := v.Validate(&schema.SendOtpRequest{Phone: "9989012345678"})

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
}

func TestSendOtpRequestAcceptsWellFormedPhone(t *testing.T) {
	v := myvalidator.NewRequestValidator(validator.New())

	if err := v.Validate(&schema.SendOtpRequest{Phone: "998901234567"}); err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
}
