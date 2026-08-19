package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type JournalBasicEntity struct {
	ID                 uint                     `json:"id"`
	Name               map[string]string        `json:"name"`
	ShortName          map[string]string        `json:"short_name"`
	Description        map[string]string        `json:"description"`
	ISSNPaper          *string                  `json:"issn"`
	ISSNOnline         *string                  `json:"issn_online"`
	DateOfEstablished  time.Time                `json:"date_of_established"`
	Website            *string                  `json:"website"`
	Address            map[string]any           `json:"address"`
	PhoneNumber        *string                  `json:"phone_number"`
	CoverImageFile     *string                  `json:"cover_image_url"`
	PublisherID        uint                     `json:"publisher_id"`
	Publisher          *PublisherEntity         `json:"publisher"`
	StudyFields        []*StudyFieldEntity      `json:"study_fields"`
	AccessType         enum.AccessType          `json:"access_type"`
	IsActive           bool                     `json:"is_active"`
	RatingCount        int64                    `json:"rating_count"`
	RatingSum          int64                    `json:"rating_sum"`
	RatingAvg          float64                  `json:"rating_avg"`
	ViewsCount         int64                    `json:"views_count"`
	OAKCertificateFile *string                  `json:"oak_certificate_file"`
	Indexes            []*JournalIndexingEntity `json:"indexes"`

	RegionID   *uint           `json:"region_id,omitempty"`
	DistrictID *uint           `json:"district_id,omitempty"`
	Region     *RegionEntity   `json:"region,omitempty"`
	District   *DistrictEntity `json:"district,omitempty"`
}

func NewJournalBasicEntity(ID uint, name map[string]string, shortName map[string]string, description map[string]string, ISSNPaper, ISSNOnline *string, dateOfEstablished time.Time, website *string, address map[string]any, phoneNumber *string, coverImageFile *string, publisherID uint, publisher *PublisherEntity, studyFields []*StudyFieldEntity, accessType enum.AccessType, isActive bool, ratingCount int64, ratingSum int64, ratingAvg float64, viewsCount int64, oakCertificateFile *string, indexes []*JournalIndexingEntity) *JournalBasicEntity {
	return &JournalBasicEntity{ID: ID, Name: name, ShortName: shortName, Description: description, ISSNPaper: ISSNPaper, ISSNOnline: ISSNOnline, DateOfEstablished: dateOfEstablished, Website: website, Address: address, PhoneNumber: phoneNumber, CoverImageFile: coverImageFile, PublisherID: publisherID, Publisher: publisher, StudyFields: studyFields, AccessType: accessType, IsActive: isActive, RatingCount: ratingCount, RatingSum: ratingSum, RatingAvg: ratingAvg, ViewsCount: viewsCount, OAKCertificateFile: oakCertificateFile, Indexes: indexes}
}
