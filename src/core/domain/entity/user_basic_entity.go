package entity

type UserBasicEntity struct {
	ID        uint   `json:"id"`
	ScienceID string `json:"science_id"`
	//Pin         string         `json:"pin"`
	FullName    string            `json:"full_name"`
	PhoneNumber string            `json:"phone_number"`
	IsAdmin     bool              `json:"is_admin"`
	Roles       []*UserRoleEntity `json:"roles"`
}

func NewUserBasicEntity(ID uint, scienceId, fullName, phoneNumber string, isAdmin bool, roles []*UserRoleEntity) *UserBasicEntity {
	return &UserBasicEntity{ID: ID, ScienceID: scienceId, FullName: fullName, PhoneNumber: phoneNumber, IsAdmin: isAdmin, Roles: roles}
}
