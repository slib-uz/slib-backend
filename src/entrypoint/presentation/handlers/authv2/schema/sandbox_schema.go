package schema

type SandboxLoginRequest struct {
	PhoneNumber string `json:"phone" validate:"required,phone_uz" example:"998901234567"`
	Otp         string `json:"otp" validate:"required"`
}
