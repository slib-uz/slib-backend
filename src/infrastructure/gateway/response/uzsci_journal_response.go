package response

import "time"

type UzSciJournalByISSNResponse struct {
	Success bool              `json:"success"`
	Data    *UzSciJournalData `json:"data"`
}

type UzSciJournalData struct {
	ID                   uint                    `json:"id"`
	Name                 string                  `json:"name"`
	ShortName            string                  `json:"short_name"`
	ISSNPrint            string                  `json:"issn_print"`
	ISSNOnline           string                  `json:"issn_online"`
	Description          string                  `json:"description"`
	Publisher            *UzSciJournalPublisher  `json:"publisher"`
	WebsiteURL           string                  `json:"website_url"`
	EstablishedDate      *time.Time              `json:"established_date"`
	CertificateFileURL   string                  `json:"certificate_file_url"`
	Email                string                  `json:"email"`
	Address              map[string]any          `json:"address"`
	PhoneNumber          string                  `json:"phone_number"`
	CoverImageURL        string                  `json:"cover_image_url"`
	PublishingPrice      float64                 `json:"publishing_price"`
	OakCertificateFileURL string                 `json:"oak_certificate_file_url"`
	IsActive             bool                    `json:"is_active"`
	AccessType           int                     `json:"access_type"`
	SubmissionAccess     int                     `json:"submission_access"`
	SocialNetworks       map[string]any          `json:"social_networks"`
	SupportLink          string                  `json:"support_link"`
}

type UzSciJournalPublisher struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	TIN         string `json:"tin"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	WebsiteURL  string `json:"website_url"`
	Address     string `json:"address"`
}
