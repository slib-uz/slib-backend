package schema

type ProfileUpdateRequest struct {
	Bio   *string `json:"bio"`
	Email string  `json:"email"`
	Photo *string `json:"photo"`
}
