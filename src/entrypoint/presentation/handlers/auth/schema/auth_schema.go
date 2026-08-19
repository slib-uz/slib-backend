package schema

// LogoutRequest — refresh token ixtiyoriy: bosqichli rollout davrida
// eski frontend uni yubormasligi mumkin.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest — refresh tokenni body orqali qabul qilish uchun.
// Query parametr 1-bosqichda orqaga moslik uchun qoladi.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse — /auth/refresh-token javobi.
//
// Javob ataylab "yassi": u response.Response ichiga o'ralmaydi. Ishlab turgan
// frontend aynan shu shaklni kutadi, bosqichli rollout esa uni buzmaslik uchun
// bor. Shakl o'zgarsa, bu tur ham, handler ham birga o'zgarishi shart.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
