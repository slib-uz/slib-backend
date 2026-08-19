package entity

type ArticleCreateEntity struct {
	Name                     map[string]string
	CoAuthorsCount           uint
	CoAuthorsIDs             []uint
	StudyFieldsIDs           []uint
	LanguageID               uint
	ContentFile              string
	TagIds                   []uint
	DOI                      *string
	JournalID                uint
	ExpertConclusionFilePath *string
	Annotation               map[string]string
	References               []string
}

func NewArticleCreateEntity(
	name map[string]string,
	coAuthorsCount uint,
	coAuthorsIDs []uint,
	studyFieldsIDs []uint,
	languageID uint,
	contentFile string,
	tagIds []uint,
	DOI *string,
	journalID uint,
	expertConclusionFilePath *string,
	annotation map[string]string,
	references []string,
) *ArticleCreateEntity {
	return &ArticleCreateEntity{
		Name:                     name,
		CoAuthorsCount:           coAuthorsCount,
		CoAuthorsIDs:             coAuthorsIDs,
		StudyFieldsIDs:           studyFieldsIDs,
		LanguageID:               languageID,
		ContentFile:              contentFile,
		TagIds:                   tagIds,
		DOI:                      DOI,
		JournalID:                journalID,
		ExpertConclusionFilePath: expertConclusionFilePath,
		Annotation:               annotation,
		References:               references,
	}
}
