package response

type Response struct {
	Status   int     `json:"status"`
	Code     *int    `json:"code,omitempty"`
	CodeName *string `json:"_code,omitempty"`
	Message  string  `json:"message,omitempty"`
	Data     any     `json:"data,omitempty"`
}

func (this *Response) Error() string {
	return this.Message
}

func NewResponse(status int, code *int, codeName *string, message string, data any) *Response {
	return &Response{Status: status, Code: code, CodeName: codeName, Message: message, Data: data}
}

func NewFailResponse(status int, message string) *Response {
	return &Response{Status: status, Message: message}
}

func NewOptionalResponse(status int, code *Code, data any, message string) *Response {
	if code != nil {
		return &Response{Status: status, Code: &code.Code, CodeName: &code.Name, Data: data, Message: message}
	}
	return &Response{Status: status, Data: data}
}

var (
	ConflictError         = NewOptionalResponse(200, CodeConflictError, nil, "")
	NotFoundError         = NewFailResponse(404, "not found")
	InvalidTokenError     = NewFailResponse(401, "invalid token")
	ExpiredTokenError     = NewFailResponse(401, "expired token")
	Empty                 = NewFailResponse(400, "response is empty")
	UnauthorizedError     = NewFailResponse(401, "unauthorized")
	PermissionDeniedError = NewFailResponse(403, "permission denied")
	InvalidCodeError      = NewFailResponse(400, "invalid code")
	ExternalServiceError  = NewFailResponse(500, "external service error")
	InvalidArgument       = NewFailResponse(400, "invalid argument")

	InvalidFileError   = NewFailResponse(400, "invalid file")
	EntityToLargeError = NewFailResponse(413, "entity too large")

	// Brute-force himoyasi (CWE-307): urinish/so'rov chegarasi oshdi.
	TooManyAttemptsError = NewFailResponse(429, "juda ko'p urinish, yangi kod so'rang")
	TooManyRequestsError = NewFailResponse(429, "juda ko'p so'rov, biroz kuting")

	// InvalidSortFieldError — so'rovda ruxsat etilmagan tartiblash maydoni
	// kelgan. 400, chunki bu mijoz xatosi; jimgina standart tartibga qaytish
	// frontend nosozligini yashirib qo'yardi.
	InvalidSortFieldError = NewFailResponse(400, "invalid sort field")

	// LogoutFailedError — tokenni denylist'ga yozib bo'lmadi. 503, chunki 200
	// qaytarish foydalanuvchiga "chiqdingiz" deb yolg'on aytadi, token esa
	// ishlashda davom etadi.
	LogoutFailedError = NewFailResponse(503, "chiqish amalga oshmadi, qaytadan urinib ko'ring")

	InsufficientBalance = NewOptionalResponse(200, CodeInsufficientBalance, nil, "")
	NotFound            = NewOptionalResponse(200, CodeNotFound, nil, "not found")
	EditionHasArticles  = NewOptionalResponse(200, CodeEditionHasArticles, nil, "")

	JournalIDHeaderRequired = NewFailResponse(400, "JournalID header is required")
	JournalIDHeaderMismatch = NewFailResponse(400, "JournalID header does not match client's journal")

	// InvalidPhoneNumberMessage va InvalidPhoneNumberError — telefon raqami
	// formati buzilganda beriladigan yagona javob. Matn ikki joyda kerak:
	// go-playground tagi (body so'rovlari) va query parametrini tekshiruvchi
	// use case. Shuning uchun u shu yerda, bitta manbada saqlanadi.
	InvalidPhoneNumberMessage = "telefon raqami noto'g'ri formatda, kutilgan format: 998XXXXXXXXX"
	InvalidPhoneNumberError   = NewFailResponse(400, InvalidPhoneNumberMessage)
)
