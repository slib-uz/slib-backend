package validator_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

// phoneFormatMessage — mijozga ko'rsatiladigan matn. Test uni ataylab
// nusxa qilib saqlaydi: konstantani import qilsak, matn tasodifan
// o'zgarganda test buni sezmay qolardi.
const phoneFormatMessage = "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"

// phoneProbe — phone_uz qoidasini tekshirish uchun lokal struktura. Haqiqiy
// schema'dan mustaqil, shuning uchun bu testlar send_otp_schema o'zgarsa ham
// o'z ma'nosini yo'qotmaydi.
type phoneProbe struct {
	Phone string `json:"phone" validate:"required,phone_uz"`
}

// plainProbe — phone_uz ishlatmaydigan struktura. phone_uz dan boshqa taglar
// xom xabarini saqlab qolishini tekshirish uchun.
type plainProbe struct {
	Title string `json:"title" validate:"required"`
}

// twoFieldProbe — telefon yonida boshqa majburiy maydon bo'lgan holat, ya'ni
// SandboxLoginRequest ning shakli. Ikkala maydon ham buzilganda foydalanuvchi
// ikkalasi haqida ham bilishi kerak.
type twoFieldProbe struct {
	Phone string `json:"phone" validate:"required,phone_uz"`
	Otp   string `json:"otp" validate:"required"`
}

// asFailResponse xatolikni *response.Response turiga keltiradi yoki testni
// muvaffaqiyatsiz deb belgilaydi.
func asFailResponse(t *testing.T, err error) *response.Response {
	t.Helper()

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	return resp
}

func newValidator() *myvalidator.RequestValidator {
	return myvalidator.NewRequestValidator(validator.New())
}

func TestValidateAcceptsWellFormedPhone(t *testing.T) {
	if err := newValidator().Validate(&phoneProbe{Phone: "998901234567"}); err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
}

func TestValidateRejectsMalformedPhone(t *testing.T) {
	cases := []struct {
		name  string
		phone string
	}{
		{"bir raqam ortiq", "9989012345678"},
		{"998 prefiksi yo'q", "901234567"},
		{"plus bilan", "+998901234567"},
		{"harfli", "abc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: tc.phone}))

			if resp.Status != 400 {
				t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
			}
			if resp.Message != phoneFormatMessage {
				t.Errorf("telefon xabari kutilgandi, %q keldi", resp.Message)
			}
		})
	}
}

// Xom go-playground xabari struct nomlarini oshkor qiladi ("Key:
// 'phoneProbe.Phone' ..."). Mijozga u yetib bormasligi kerak.
func TestValidateDoesNotLeakStructNameForPhone(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: "abc"}))

	if strings.Contains(resp.Message, "phoneProbe") || strings.Contains(resp.Message, "Key:") {
		t.Errorf("xabar ichki tafsilotni oshkor qildi: %q", resp.Message)
	}
}

// Bo'sh maydon phone_uz ni ham buzadi, lekin sabab formatda emas — maydon
// umuman to'ldirilmagan. Foydalanuvchi "required" haqida eshitishi kerak.
func TestValidateEmptyPhoneReportsRequiredNotFormat(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&phoneProbe{Phone: ""}))

	if resp.Message == phoneFormatMessage {
		t.Fatal("bo'sh maydon uchun 'required' xabari kutilgandi, format xabari keldi")
	}
	if !strings.Contains(resp.Message, "required") {
		t.Errorf("'required' tagi haqida xabar kutilgandi, %q keldi", resp.Message)
	}
}

// phone_uz dan boshqa taglar hozirgi xulqini saqlaydi — bu o'zgarish mavjud
// endpointlarning javobiga tegmasligi kerak.
func TestValidateLeavesOtherTagMessagesUnchanged(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&plainProbe{Title: ""}))

	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
	if !strings.Contains(resp.Message, "required") {
		t.Errorf("xom validator xabari kutilgandi, %q keldi", resp.Message)
	}
}

// Telefon buzuq va otp bo'sh bo'lsa, ikkala xato ham mavjud. Faqat telefon
// xabarini qaytarish otp haqidagi ma'lumotni yashiradi — aynan shu niqoblanish
// ushlanishi kerak.
func TestValidateDoesNotMaskSiblingFieldError(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&twoFieldProbe{Phone: "abc", Otp: ""}))

	if resp.Message == phoneFormatMessage {
		t.Fatalf("xabar faqat telefon haqida edi, otp xatosi yashirildi: %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "otp") && !strings.Contains(resp.Message, "Otp") {
		t.Errorf("xabarda otp maydoni tilga olinishi kutilgandi, %q keldi", resp.Message)
	}
}

// Telefon buzuq, lekin otp to'ldirilgan bo'lsa — yagona xato phone_uz, demak
// tuzatish telefon xabarini butunlay o'chirib yubormagan bo'lishi kerak.
func TestValidateKeepsPhoneMessageWhenItIsTheOnlyError(t *testing.T) {
	resp := asFailResponse(t, newValidator().Validate(&twoFieldProbe{Phone: "abc", Otp: "111111"}))

	if resp.Message != phoneFormatMessage {
		t.Errorf("telefon xabari kutilgandi, %q keldi", resp.Message)
	}
}
