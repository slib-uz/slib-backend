package schema

type SendOtpRequest struct {
	Phone string `json:"phone" validate:"required,phone_uz" example:"998901234567"`
}

type SendOtpResponse struct {
	SessionID string `json:"session_id"`
}
