package entity

type PublisherAntiplagStatsEntity struct {
	PublisherID   uint   `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	Success       int    `json:"success"`
	Failed        int    `json:"failed"`
}

func NewPublisherAntiplagStatsEntity(publisherID uint, publisherName string, success int, failed int) *PublisherAntiplagStatsEntity {
	return &PublisherAntiplagStatsEntity{PublisherID: publisherID, PublisherName: publisherName, Success: success, Failed: failed}
}
