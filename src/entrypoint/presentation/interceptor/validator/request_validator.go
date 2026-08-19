package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/utils"
)

// phoneTag — O'zbekiston telefon raqami formatini tekshiruvchi custom qoida.
// Schema'da `validate:"required,phone_uz"` ko'rinishida ishlatiladi; tartib
// muhim, chunki go-playground birinchi buzilgan tagda to'xtaydi va bo'sh
// maydon uchun format emas, "required" xabari chiqishi kerak.
const phoneTag = "phone_uz"

type RequestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator(validator *validator.Validate) *RequestValidator {
	// RegisterValidation faqat tag nomi bo'sh bo'lganda xato qaytaradi, bu
	// yerda esa konstanta. Qoida shu yerda ro'yxatdan o'tadi (NewEcho da
	// emas), toki validator qurilgan har bir joyda — testlarni ham qo'shib —
	// u mavjud bo'lsin.
	_ = validator.RegisterValidation(phoneTag, isValidPhone)

	return &RequestValidator{validator: validator}
}

// isValidPhone — phoneTag qoidasining tanasi. Format qoidasining o'zi
// utils da, bitta joyda saqlanadi.
func isValidPhone(fl validator.FieldLevel) bool {
	return utils.IsValidPhoneNumber(fl.Field().String())
}

func (this *RequestValidator) Validate(i interface{}) error {
	if err := this.validator.Struct(i); err != nil {
		return response.NewFailResponse(400, validationMessage(err))
	}
	if v, ok := i.(Validatable); ok {
		if ok, err := v.Validate(); !ok && err != nil {
			return response.NewFailResponse(400, err.Error())
		}
	}
	return nil
}

// validationMessage phoneTag buzilgan bo'lsa tushunarli matn qaytaradi, aks
// holda validator xatosini o'zgarishsiz uzatadi — qolgan endpointlarning
// javobi o'zgarmasligi uchun. Telefon xabari faqat phone_uz yagona buzilgan
// qoida bo'lganda beriladi: bitta xato bo'lsa sabab aniq, lekin bir nechta
// maydon buzilganda faqat telefon haqida aytish qo'shni maydon xatosini
// (masalan, bo'sh otp) mijozdan yashiradi — shuning uchun bunday holda xom
// go-playground xabari qaytadi, u hammasini sanab beradi.
func validationMessage(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) && len(errs) == 1 && errs[0].Tag() == phoneTag {
		return response.InvalidPhoneNumberMessage
	}
	return err.Error()
}
