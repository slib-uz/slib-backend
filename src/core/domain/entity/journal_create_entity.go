package entity

import (
	enum "slib.uz/src/core/domain/entity/enum"
)

type JournalCreateEntity struct {
	Name                     map[string]string        `json:"name"`
	ShortName                map[string]string        `json:"short_name"`
	ISSNPaper                *string                  `json:"issn_paper"`  // ISSNPaper for paper version
	ISSNOnline               *string                  `json:"issn_online"` // ISSNPaper for online version
	Description              map[string]string        `json:"description"`
	ArticlePublishConditions map[string]string        `json:"article_publish_conditions"`
	Rule                     map[string]string        `json:"rule"`
	DateOfEstablished        string                   `json:"date_of_established"`
	Website                  *string                  `json:"website"`
	Indexes                  []*JournalIndexingEntity `json:"indexing"`
	CertificateFile          *string                  `json:"certificate_url"` // Guvohnomasi
	Email                    *string                  `json:"email"`
	Address                  map[string]any           `json:"address"`
	PhoneNumber              *string                  `json:"phone_number"`
	CoverImageFile           *string                  `json:"cover_image_url"`
	PublisherID              uint                     `json:"publisher_id"`
	StudyFieldIDs            []uint                   `json:"study_field_id"`
	PublishingPrice          *int64                   `json:"publishing_price"`    // pricing for publishing an article
	SellingPrice             *int64                   `json:"selling_price"`       // pricing for selling a journal
	AccessType               enum.AccessType          `json:"access_type"`         // Access type of the journal
	OAKCertificateFile       *string                  `json:"oak_certificate_url"` // OAK guvohnomasi
	LanguageIDs              []uint                   `json:"language_ids"`
	PeerReviewMethod         enum.PeerReviewMethod    `json:"peer_review_method"` // Method of peer review

	SubmissionAccess enum.AccessType `json:"submission_access"` // Access type for article submission
	CommentAccess    enum.AccessType `json:"comment_access"`    // Access type for comments
	SocialNetworks   map[string]any  `json:"social_networks"`
	SupportLink      *string         `json:"support_link"`

	RegionID   *uint `json:"region_id,omitempty"`
	DistrictID *uint `json:"district_id,omitempty"`
}

func NewJournalCreateEntity(
	name map[string]string,
	shortName map[string]string,
	ISSNPaper *string,
	ISSNOnline *string,
	description map[string]string,
	articlePublishConditions map[string]string,
	rule map[string]string,
	dateOfEstablished string,
	website *string,
	indexes []*JournalIndexingEntity,
	certificateFile *string,
	email *string,
	address map[string]any,
	phoneNumber *string,
	coverImageFile *string,
	publisherID uint,
	studyFieldIDs []uint,
	publishingPrice *int64,
	sellingPrice *int64,
	accessType enum.AccessType,
	OAKCertificateFile *string,
	languageIDs []uint,
	peerReviewMethod enum.PeerReviewMethod,
	submissionAccess enum.AccessType,
	commentAccess enum.AccessType,
	socialNetworks map[string]any,
	supportLink *string,
) *JournalCreateEntity {
	return &JournalCreateEntity{
		Name:                     name,
		ShortName:                shortName,
		ISSNPaper:                ISSNPaper,
		ISSNOnline:               ISSNOnline,
		Description:              description,
		ArticlePublishConditions: articlePublishConditions,
		Rule:                     rule,
		DateOfEstablished:        dateOfEstablished,
		Website:                  website,
		Indexes:                  indexes,
		CertificateFile:          certificateFile,
		Email:                    email,
		Address:                  address,
		PhoneNumber:              phoneNumber,
		CoverImageFile:           coverImageFile,
		PublisherID:              publisherID,
		StudyFieldIDs:            studyFieldIDs,
		PublishingPrice:          publishingPrice,
		SellingPrice:             sellingPrice,
		AccessType:               accessType,
		OAKCertificateFile:       OAKCertificateFile,
		LanguageIDs:              languageIDs,
		PeerReviewMethod:         peerReviewMethod,
		SubmissionAccess:         submissionAccess,
		CommentAccess:            commentAccess,
		SocialNetworks:           socialNetworks,
		SupportLink:              supportLink}
}
