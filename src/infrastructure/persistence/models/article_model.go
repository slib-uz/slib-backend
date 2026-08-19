package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type ArticleModel struct {
	gorm.Model

	Name            datatypes.JSON
	PublicationDate datatypes.Date
	CoAuthorsCount  int                `gorm:"not null"` // total co-authors count
	CoAuthors       []*AuthorModel     `gorm:"many2many:article_co_authors;constraint:OnDelete:CASCADE;"`
	AccessType      enum.AccessType    `gorm:"not null"`
	StudyFields     []*StudyFieldModel `gorm:"many2many:article_study_fields;constraint:OnDelete:CASCADE;"`
	LanguageID      uint               `gorm:"not null"`
	Language        *LanguageModel     `gorm:"foreignKey:LanguageID;references:ID; constraint:OnDelete:SET NULL"`
	Annotation      datatypes.JSON
	// ContentFileID   uuid.UUID          `gorm:"type:uuid;"`
	ContentFilePath string
	Tags            []*TagModel `gorm:"many2many:article_tags;constraint:OnDelete:CASCADE;"`

	DOI *string `gorm:"size:255;uniqueIndex"`
	ROI *string `gorm:"size:255;uniqueIndex"`

	JournalID uint
	Journal   *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL"`

	ExpertConclusionFile *string
	IsPublished          bool `gorm:"default:false;not null"`

	ViewsCount int64 `gorm:"default:0;not null"` // total views count

	RatingCount int64 `gorm:"default:0;not null"` // total rating count
	RatingSum   int64 `gorm:"default:0;not null"` // total rating sum

	Reports []ReportModel `gorm:"polymorphic:Target;"`

	ArticleAuthorAffiliations []*ArticleAuthorAffiliationModel `gorm:"foreignKey:ArticleID;"`

	ExternalID *uint `gorm:"uniqueIndex;"`
	IsSent     bool  `gorm:"default:false;not null"`

	EditionID *uint
	Edition   *EditionModel `gorm:"foreignKey:EditionID;references:ID;constraint:OnDelete:SET NULL"`

	Pages      string `gorm:"size:255"`
	WebAddress string `gorm:"size:512"`

	References []*ReferenceModel `gorm:"-"`

	// Old SLIB database dagi mualliflar. Bular faqat ko'rsatish uchun.
	UnconfirmedAuthors datatypes.JSON `gorm:"type:jsonb"`
}

func NewArticleModel(id uint, name datatypes.JSON, coAuthorsCount int, coAuthors []*AuthorModel, accessType enum.AccessType, studyFields []*StudyFieldModel, languageID uint, annotation datatypes.JSON, file string, tags []*TagModel, DOI *string, journalID uint, expertConclusionFile *string) *ArticleModel {
	return &ArticleModel{
		Model:                gorm.Model{ID: id},
		Name:                 name,
		CoAuthorsCount:       coAuthorsCount,
		CoAuthors:            coAuthors,
		AccessType:           accessType,
		StudyFields:          studyFields,
		LanguageID:           languageID,
		Annotation:           annotation,
		ContentFilePath:      file,
		Tags:                 tags,
		DOI:                  DOI,
		JournalID:            journalID,
		ExpertConclusionFile: expertConclusionFile,
	}
}

func (ArticleModel) TableName() string {
	return "articles"
}
