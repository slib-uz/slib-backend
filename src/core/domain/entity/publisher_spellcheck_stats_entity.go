package entity

type PublisherSpellcheckStatsEntity struct {
	PublisherID   uint   `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	Success       int    `json:"success"`
	Failed        int    `json:"failed"`
}

func NewPublisherSpellcheckStatsEntity(publisherID uint, publisherName string, success int, failed int) *PublisherSpellcheckStatsEntity {
	return &PublisherSpellcheckStatsEntity{PublisherID: publisherID, PublisherName: publisherName, Success: success, Failed: failed}
}
