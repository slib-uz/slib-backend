package response

type AntiPlagResultResponse struct {
	ExternalID        uint    `json:"id"`
	Status            int     `json:"status"`
	StatusDisplay     string  `json:"status_display"`
	PlagiarismPercent float64 `json:"plagiarism_percent"`
	LegalPercent      float64 `json:"legal_percent"`
	SelfCitePercent   float64 `json:"self_cite_percent"`
	UnknownPercent    float64 `json:"unknown_percent"`
	ShortReportURL    string  `json:"short_report_url"`
	FullReportURL     string  `json:"full_report_url"`
	ExternalCreatedAt string  `json:"created_at"`
}

type AntiPlagBalanceResponse struct {
	Balance int `json:"balance"`
}
