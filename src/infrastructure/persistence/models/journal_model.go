package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalModel struct {
	gorm.Model

	Name                     datatypes.JSON
	ShortName                datatypes.JSON
	ISSNPaper                *string `gorm:"size:255;index"`
	ISSNOnline               *string `gorm:"size:255;index"`
	Description              datatypes.JSON
	Rule                     datatypes.JSON
	ArticlePublishConditions datatypes.JSON
	EstablishedDate          time.Time
	Website                  *string                 `gorm:"size:128"`
	Indexes                  []*JournalIndexingModel `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE;"`
	CertificateFile          *string                 `gorm:"size:500"` // Guvohnomasi
	Email                    *string                 `gorm:"size:255"`
	Address                  datatypes.JSON
	PhoneNumber              *string `gorm:"size:255"`
	CoverImageFile           *string `gorm:"size:500"`

	PublisherID uint            `gorm:"index;not null;"`
	Publisher   *PublisherModel `gorm:"foreignKey:PublisherID;constraint:OnDelete:SET NULL;"`

	StudyFields []*StudyFieldModel `gorm:"many2many:journal_many2many_study_fields;"`

	PublishingPrice *int64
	SellingPrice    *int64
	AccessType      enum.AccessType

	OAKCertificateFile *string `gorm:"size:500"` // OAK guvohnomasi (sertifikati)

	Languages        []*LanguageModel      `gorm:"many2many:journal_languages;"`
	PeerReviewMethod enum.PeerReviewMethod `gorm:"not null"`

	IsActive   bool `gorm:"default:false;not null"`
	IsAccepted bool `gorm:"default:false;not null"`

	SubmissionAccess enum.AccessType `gorm:"not null;default:10"`
	CommentAccess    enum.AccessType `gorm:"not null;default:10"`

	SocialNetworks datatypes.JSON `gorm:"type:jsonb;"`
	SupportLink    *string        `gorm:"size:255;"`

	Reviewers []*ReviewerModel `gorm:"many2many:journal_reviewers;constraint:OnDelete:CASCADE"`
	Reports   []ReportModel    `gorm:"polymorphic:Target;"`

	RatingCount int64 `gorm:"default:0;not null"` // total rating count
	RatingSum   int64 `gorm:"default:0;not null"` // total rating sum

	ViewsCount int64 `gorm:"default:0;not null"` // total views count

	RegionID   *uint          `gorm:"index"`
	Region     *RegionModel   `gorm:"foreignKey:RegionID;constraint:OnDelete:SET NULL;"`
	DistrictID *uint          `gorm:"index"`
	District   *DistrictModel `gorm:"foreignKey:DistrictID;constraint:OnDelete:SET NULL;"`
}

func NewJournalModel(name datatypes.JSON, shortName datatypes.JSON, ISSNPaper *string, ISSNOnline *string, description datatypes.JSON, rule datatypes.JSON, articlePublishConditions datatypes.JSON, establishedDate time.Time, website *string, certificateFile *string, email *string, address datatypes.JSON, phoneNumber *string, coverImageFile *string, publisherID uint, studyFields []*StudyFieldModel, publishingPrice, sellingPrice *int64, accessType enum.AccessType, oakCertificate *string, peerReviewMethod enum.PeerReviewMethod, submissionAccess enum.AccessType, CommentAccess enum.AccessType, socialNetworks datatypes.JSON, supportLink *string) *JournalModel {
	return &JournalModel{Name: name, ShortName: shortName, ISSNPaper: ISSNPaper, ISSNOnline: ISSNOnline, Description: description, Rule: rule, ArticlePublishConditions: articlePublishConditions, EstablishedDate: establishedDate, Website: website, CertificateFile: certificateFile, Email: email, Address: address, PhoneNumber: phoneNumber, CoverImageFile: coverImageFile, PublisherID: publisherID, StudyFields: studyFields, PublishingPrice: publishingPrice, SellingPrice: sellingPrice, AccessType: accessType, OAKCertificateFile: oakCertificate, PeerReviewMethod: peerReviewMethod, SubmissionAccess: submissionAccess, CommentAccess: CommentAccess, SocialNetworks: socialNetworks, SupportLink: supportLink}
}

func (*JournalModel) TableName() string {
	return "journals"
}

func (this *JournalModel) BeforeCreate(tx *gorm.DB) (err error) {
	if this.IsAccepted {
		var count int64
		if this.ISSNPaper != nil {
			tx.Model(this).Where("issn_paper = ? and is_accepted = ?", this.ISSNPaper, true).Count(&count)
			if count > 0 {
				return fmt.Errorf("journal with this ISSNPaper already exists")
			}
		}
		if this.ISSNOnline != nil {
			tx.Model(this).Where("issn_online = ? and is_accepted = ?", this.ISSNOnline, true).Count(&count)
			if count > 0 {
				return fmt.Errorf("journal with this ISSNOnline already exists")
			}
		}
		if this.ISSNPaper == nil && this.ISSNOnline == nil {
			return fmt.Errorf("at least one of ISSNPaper or ISSNOnline must be provided")
		}
	}
	return nil
}
