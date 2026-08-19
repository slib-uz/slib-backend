package entity

type PartnerEntity struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	LogoPath string `json:"logo"`
	Link     string `json:"link"`
}

func NewPartnerEntity(id uint, title string, logoPath string, link string) *PartnerEntity {
	return &PartnerEntity{ID: id, Title: title, LogoPath: logoPath, Link: link}
}
