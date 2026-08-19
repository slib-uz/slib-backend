package response

type AiDetectResultResponse struct {
	ExternalID    uint    `json:"id"`
	WordsCount    int     `json:"words_count"`
	Status        int     `json:"status"`
	StatusDisplay string  `json:"status_display"`
	HumanPercent  float64 `json:"human_percent"`
	WarnPercent   float64 `json:"warn_percent"`
	AiPercent     float64 `json:"ai_percent"`
	ReportURL     string  `json:"report_url"`
	CreatedAt     string  `json:"created_at"`
}
