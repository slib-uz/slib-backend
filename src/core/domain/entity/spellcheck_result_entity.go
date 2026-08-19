package entity

import "time"

type SpellCheckResultEntity struct {
	ID            uint                       `json:"id"`
	ReviewStageID uint                       `json:"review_stage_id"`
	ReviewStage   *ReviewStageResponseEntity `json:"review_stage"`
	ApplicationID uint                       `json:"application_id" `
	Application   *ApplicationBasicEntity    `json:"application"`
	JournalID     uint                       `json:"journal_id"`
	Journal       *JournalEntity             `json:"journal"`
	File          string                     `json:"file"`
	ResultFile    *string                    `json:"result_file"`
	Status        int                        `json:"status"`
	SubmitterID   uint                       `json:"submitter_id"`
	Submitter     *UserBasicEntity           `json:"submitter"`
	ResultTime    *time.Time                 `json:"result_time"`
}

func NewSpellCheckResultEntity(ID uint, reviewStageID uint, reviewStage *ReviewStageResponseEntity, applicationID uint, application *ApplicationBasicEntity, journalID uint, journal *JournalEntity, file string, resultFile *string, status int, submitterID uint, submitter *UserBasicEntity, resultTime *time.Time) *SpellCheckResultEntity {
	return &SpellCheckResultEntity{ID: ID, ReviewStageID: reviewStageID, ReviewStage: reviewStage, ApplicationID: applicationID, Application: application, JournalID: journalID, Journal: journal, File: file, ResultFile: resultFile, Status: status, SubmitterID: submitterID, Submitter: submitter, ResultTime: resultTime}
}
