package schema

import "slib.uz/src/core/domain/entity/enum"

type LoginRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Code      string `json:"code" validate:"required"`
	ScienceID string `json:"science_id"`
	FullName  string `json:"full_name"`
}

type LoginResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	Scope        enum.AuthScope `json:"scope"`
}
