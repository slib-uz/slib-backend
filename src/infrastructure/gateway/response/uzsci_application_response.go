package response

import "time"

type UzSciCreateApplicationResponse struct {
	Success bool                  `json:"success"`
	Data    *UzSciApplicationData `json:"data"`
}

type UzSciApplicationData struct {
	ID             uint                         `json:"id"`
	RatingPeriodID uint                         `json:"rating_period_id"`
	RatingPeriod   *UzSciRatingPeriodResponse   `json:"rating_period"`
	JournalID      uint                         `json:"journal_id"`
	Journal        *UzSciJournalData            `json:"journal"`
	Status         string                       `json:"status"`
	Answers        []UzSciApplicationAnswerData `json:"answers"`
	CreatedAt      time.Time                    `json:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at"`
}

type UzSciApplicationAnswerData struct {
	ID            uint               `json:"id"`
	ApplicationID uint               `json:"application_id"`
	FormID        uint               `json:"form_id"`
	Form          *UzSciFormResponse `json:"form"`
	Value         string             `json:"value"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type UzSciCreateApplicationPayload struct {
	JournalID uint                               `json:"journal_id"`
	Answers   []UzSciCreateApplicationAnswerItem `json:"answers"`
}

type UzSciCreateApplicationAnswerItem struct {
	FormID uint   `json:"form_id"`
	Value  string `json:"value"`
}
