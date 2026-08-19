package entity

import (
	"slib.uz/src/core/utils"
)

type RoiPublishGatewayEntity struct {
	Name             map[string]string
	PublicationDate  utils.DateOnlyType
	CoAuthorsCount   uint
	CoAuthors        []*AuthorEntity
	AccessType       int
	StudyFieldsCodes []uint
	LanguageID       uint
	File             string
	DOI              *string
	JournalID        uint
	Citations        []string
}

func NewRoiPublishGatewayEntity(
	name map[string]string,
	publicationDate utils.DateOnlyType,
	CoAuthorsCount uint,
	CoAuthors []*AuthorEntity,
	AccessType int,
	StudyFieldsCodes []uint,
	LanguageID uint,
	File string,
	DOI *string,
	JournalID uint,
	Citations []string,
) *RoiPublishGatewayEntity {
	return &RoiPublishGatewayEntity{
		Name:             name,
		PublicationDate:  publicationDate,
		CoAuthorsCount:   CoAuthorsCount,
		CoAuthors:        CoAuthors,
		AccessType:       AccessType,
		StudyFieldsCodes: StudyFieldsCodes,
		LanguageID:       LanguageID,
		File:             File,
		DOI:              DOI,
		JournalID:        JournalID,
		Citations:        Citations,
	}
}
