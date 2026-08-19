package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type DoiSettingRequest struct {
	JournalName string `json:"journal_name" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	DOIPrefix   string `json:"doi_prefix" validate:"required"`
	DOISuffix   string `json:"doi_suffix" validate:"required"`
}

func (this *DoiSettingRequest) ToEntity(journalID uint) *entity.JournalDoiSettingEntity {
	return entity.NewJournalDoiSettingEntity(0, journalID, this.JournalName, this.Username, this.Password, this.DOIPrefix, this.DOISuffix, time.Time{}, time.Time{})
}

type CrossRefCheckAuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DoiSettingResponse struct {
	ID          uint   `json:"id"`
	JournalID   uint   `json:"journal_id"`
	JournalName string `json:"journal_name"`
	Username    string `json:"username"`
	DOIPrefix   string `json:"doi_prefix"`
	DOISuffix   string `json:"doi_suffix"`
}

func NewDoiSettingResponse(e *entity.JournalDoiSettingEntity) *DoiSettingResponse {
	return &DoiSettingResponse{
		ID:          e.ID,
		JournalID:   e.JournalID,
		JournalName: e.JournalName,
		Username:    e.Username,
		DOIPrefix:   e.DOIPrefix,
		DOISuffix:   e.DOISuffix,
	}
}
