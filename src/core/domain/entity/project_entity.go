package entity

type ProjectEntity struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	LogoPath string `json:"logo"`
	Link     string `json:"link"`
}

func NewProjectEntity(id uint, title string, logoPath string, link string) *ProjectEntity {
	return &ProjectEntity{ID: id, Title: title, LogoPath: logoPath, Link: link}
}
