package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type JournalStatisticV2Entity struct {
	JournalID               uint              `json:"journal_id"`
	JournalName             map[string]string `json:"journal_name"`
	PublisherID             uint              `json:"publisher_id"`
	PublisherName           string            `json:"publisher_name"`
	InstitutionID           uint              `json:"institution_id"`
	InstitutionName         string            `json:"institution_name"`
	IssnOnline              string            `json:"issn_online"`
	IssnPaper               string            `json:"issn_paper"`
	JournalRegionID         *uint             `json:"region_id"`
	JournalRegionName       *string           `json:"region_name"`
	JournalAccessType       enum.AccessType   `json:"access_type"`
	JournalSubmissionAccess enum.AccessType   `json:"submission_access"`
	AdminCount              int               `json:"admin_count"`
	EditionCount            int               `json:"edition_count"`
	ArticleCount            int               `json:"article_count"`
	ArticleApplicationCount int               `json:"article_application_count"`
	CoAuthorCount           int               `json:"co_author_count"`
	CompletionPercent       int               `json:"completion_percent"`
}

type JournalStatisticV2DetailEntity struct {
	ArticleCount                      int `json:"article_count"`
	ArticleApplicationCount           int `json:"article_application_count"`
	CoAuthorCount                     int `json:"coauthor_count"`
	EditionCount                      int `json:"edition_count"`
	Last30DaysArticleApplicationCount int `json:"last_30_days_article_application_count"`
	Last30DaysArticleCount            int `json:"last_30_days_article_count"`
}
