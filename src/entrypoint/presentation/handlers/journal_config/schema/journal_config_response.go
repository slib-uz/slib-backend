package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type JournalConfigResponse struct {
	ID         uint             `json:"id"`
	JournalID  uint             `json:"journal_id"`
	CreatorID  uint             `json:"creator_id"`
	WebsiteURL string           `json:"website_url"`
	Conf       map[string]any   `json:"conf"`
	IsActive   bool             `json:"is_active"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Journal    *JournalResponse `json:"journal,omitempty"`
	Creator    *CreatorResponse `json:"creator,omitempty"`
}

type JournalResponse struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
}

type CreatorResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}

func FromJournalConfigEntity(e *entity.JournalConfigEntity) *JournalConfigResponse {
	if e == nil {
		return nil
	}

	resp := &JournalConfigResponse{
		ID:         e.ID,
		JournalID:  e.JournalID,
		CreatorID:  e.CreatorID,
		WebsiteURL: e.WebsiteURL,
		Conf:       e.Conf,
		IsActive:   e.IsActive,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}

	if e.Journal != nil {
		resp.Journal = &JournalResponse{
			ID:   e.Journal.ID,
			Name: e.Journal.Name,
		}
	}

	if e.Creator != nil {
		resp.Creator = &CreatorResponse{
			ID:       e.Creator.ID,
			FullName: e.Creator.FullName,
		}
	}

	return resp
}
