package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type JournalEntity struct {
	ID                       uint                     `json:"id"`
	Name                     map[string]string        `json:"name"`
	ShortName                map[string]string        `json:"short_name"`
	ISSNPaper                *string                  `json:"issn"`
	ISSNOnline               *string                  `json:"issn_online"` // ISSNPaper for online version
	Description              map[string]string        `json:"description"`
	ArticlePublishConditions map[string]string        `json:"article_publish_conditions"`
	Rule                     map[string]string        `json:"rule"`
	DateOfEstablished        time.Time                `json:"date_of_established"`
	Website                  *string                  `json:"website"`
	Indexes                  []*JournalIndexingEntity `json:"indexing"`
	CertificateFile          *string                  `json:"certificate_url"` // Guvohnomasi
	Email                    *string                  `json:"email"`
	Address                  map[string]any           `json:"address"`
	PhoneNumber              *string                  `json:"phone_number"`
	CoverImageFile           *string                  `json:"cover_image_url"`
	PublisherID              uint                     `json:"publisher_id"`
	Publisher                *PublisherEntity         `json:"publisher"`
	StudyFields              []*StudyFieldEntity      `json:"study_fields"`
	PublishingPrice          *int64                   `json:"publishing_price"`
	SellingPrice             *int64                   `json:"selling_price"` // pricing for selling a article
	AccessType               enum.AccessType          `json:"access_type"`   // Access type of the journal, e.g., open access, subscription-based
	OAKCertificateFile       *string                  `json:"oak_certificate_url"`
	Languages                []*LanguageEntity        `json:"languages"`
	PeerReviewMethod         enum.PeerReviewMethod    `json:"peer_review_method"` // Method of peer review
	IsActive                 bool                     `json:"is_active"`
	IsAccepted               bool                     `json:"is_accepted"`
	RatingCount              int64                    `json:"rating_count"`
	RatingSum                int64                    `json:"rating_sum"`
	RatingAvg                float64                  `json:"rating_avg"`
	ViewsCount               int64                    `json:"views_count"`
	ApplicationID            uint                     `json:"application_id"`

	SubmissionAccess enum.AccessType `json:"submission_access"` // Access type for article submission
	CommentAccess    enum.AccessType `json:"comment_access"`    // Access type for comments
	SocialNetworks   map[string]any  `json:"social_networks"`
	SupportLink      *string         `json:"support_link"`

	RegionID   *uint           `json:"region_id,omitempty"`
	DistrictID *uint           `json:"district_id,omitempty"`
	Region     *RegionEntity   `json:"region,omitempty"`
	District   *DistrictEntity `json:"district,omitempty"`
}

func NewJournalEntity(ID uint, name map[string]string, shortName map[string]string, ISSNPaper *string, ISSNOnline *string, description map[string]string, articlePublishConditions map[string]string, rule map[string]string, dateOfEstablished time.Time, website *string, indexes []*JournalIndexingEntity, certificateFile *string, email *string, address map[string]any, phoneNumber *string, coverImageFile *string, publisherID uint, publisher *PublisherEntity, studyFields []*StudyFieldEntity, publishingPrice *int64, sellingPrice *int64, accessType enum.AccessType, OAKCertificateFile *string, languages []*LanguageEntity, peerReviewMethod enum.PeerReviewMethod, isActive bool, isAccepted bool, ratingCount int64, ratingSum int64, ratingAvg float64, viewsCount int64, applicationID uint, submissionAccess enum.AccessType, commentAccess enum.AccessType, socialNetworks map[string]any, supportLink *string) *JournalEntity {
	return &JournalEntity{ID: ID, Name: name, ShortName: shortName, ISSNPaper: ISSNPaper, ISSNOnline: ISSNOnline, Description: description, ArticlePublishConditions: articlePublishConditions, Rule: rule, DateOfEstablished: dateOfEstablished, Website: website, Indexes: indexes, CertificateFile: certificateFile, Email: email, Address: address, PhoneNumber: phoneNumber, CoverImageFile: coverImageFile, PublisherID: publisherID, Publisher: publisher, StudyFields: studyFields, PublishingPrice: publishingPrice, SellingPrice: sellingPrice, AccessType: accessType, OAKCertificateFile: OAKCertificateFile, Languages: languages, PeerReviewMethod: peerReviewMethod, IsActive: isActive, IsAccepted: isAccepted, RatingCount: ratingCount, RatingSum: ratingSum, RatingAvg: ratingAvg, ViewsCount: viewsCount, ApplicationID: applicationID, SubmissionAccess: submissionAccess, CommentAccess: commentAccess, SocialNetworks: socialNetworks, SupportLink: supportLink}
}
