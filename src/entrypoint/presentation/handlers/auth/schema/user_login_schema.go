package schema

type UserLoginRequest struct {
	Code string `json:"code" validate:"required"`
}
