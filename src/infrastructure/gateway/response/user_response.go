package response

type UserResponse struct {
	Status int              `json:"status"`
	Data   UserDataResponse `json:"data"`
}

type UserDataResponse struct {
	ScienceID   string `json:"science_id"`
	Pin         string `json:"pin"`
	Gender      uint   `json:"gd"`
	MiddleName  string `json:"middle_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"sur_name"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	BirthDate   string `json:"birth_date"`
}
