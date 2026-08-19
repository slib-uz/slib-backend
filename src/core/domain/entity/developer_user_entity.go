package entity

type DeveloperUserEntity struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsActive    bool   `json:"is_active"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperUser bool   `json:"is_superuser"`
}
