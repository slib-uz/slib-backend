package entity

type UnconfirmedAuthorEntity struct {
	Position         string `json:"position"`
	LastName         string `json:"last_name"`
	FirstName        string `json:"first_name"`
	MiddleName       string `json:"middle_name"`
	OrganizationName string `json:"organization_name"`
}
