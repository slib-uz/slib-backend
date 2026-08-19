package entity

import (
	"encoding/json"
)

type JournalStatisticEntity struct {
	JournalID              uint                         `json:"journal_id"`
	JournalName            string                       `json:"journal_name"`
	PublisherID            uint                         `json:"publisher_id"`
	PublisherName          string                       `json:"publisher_name"`
	IssnOnline             string                       `json:"issn_online"`
	IssnPaper              string                       `json:"issn_paper"`
	CompletionPercent      int                          `json:"completion_percent"`
	HasStudyFields         bool                         `json:"has_study_fields"`
	ArticleCount           int64                        `json:"article_count"`
	JournalPhoneNumber     string                       `json:"journal_phone_number"`
	JournalTelegramContact string                       `json:"journal_telegram_contact"`
	JournalChiefEditor     *JournalStatisticUserEntity  `json:"journal_chief_editor"`
	JournalSecretaries     []JournalStatisticUserEntity `json:"journal_secretaries"`
}

type JournalCompletionResultEntity struct {
	CompletionPercent int      `json:"completion_percent"`
	IncompleteFields  []string `json:"incomplete_fields"`
}

type JournalCompletionStatsEntity struct {
	JournalID             uint `json:"id"`
	HasName               bool `json:"has_name"`
	HasShortName          bool `json:"has_short_name"`
	HasJournalLanguages   bool `json:"has_journal_languages"`
	HasEstablishedDate    bool `json:"has_established_date"`
	HasPhoneNumber        bool `json:"has_phone_number"`
	HasEmail              bool `json:"has_email"`
	HasWebsite            bool `json:"has_website"`
	HasTelegramSupport    bool `json:"has_telegram_support"`
	HasIssnOnline         bool `json:"has_issn_online"`
	HasIssnPaper          bool `json:"has_issn_paper"`
	HasSocialNetworks     bool `json:"has_social_networks"`
	HasStudyFields        bool `json:"has_study_fields"`
	HasImageFile          bool `json:"has_image_file"`
	HasCertificateFile    bool `json:"has_certificate_file"`
	HasOakCertificateFile bool `json:"has_oak_certificate_file"`
	HasAddress            bool `json:"has_address"`
	HasJournalIndexing    bool `json:"has_journal_indexing"`
	HasDescription        bool `json:"has_description"`
	HasRule               bool `json:"has_rule"`
	HasArticleConditions  bool `json:"has_article_conditions"`
	HasRegionID           bool `json:"has_region_id"`
}

type JournalStatisticUserEntity struct {
	UserID      uint   `json:"user_id"`
	FullName    string `json:"full_name"`
	ScienceID   string `json:"science_id"`
	PhoneNumber string `json:"phone_number"`
}

func UnmarshalJournalStatisticUserEntity(data string) (*JournalStatisticUserEntity, error) {
	if data == "" || data == "[]" || data == "null" {
		return nil, nil
	}
	var result JournalStatisticUserEntity
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func UnmarshalJournalStatisticUserEntities(data string) ([]JournalStatisticUserEntity, error) {
	if data == "" || data == "[]" || data == "null" {
		return []JournalStatisticUserEntity{}, nil
	}
	var result []JournalStatisticUserEntity
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	if result == nil {
		return []JournalStatisticUserEntity{}, nil
	}
	return result, nil
}
