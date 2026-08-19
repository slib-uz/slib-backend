package entity

type ArticleUpdateEntity struct {
	Name           *map[string]string
	CoAuthorsCount *int
	CoAuthorsIDs   *[]uint
	StudyFieldsIDs *[]uint
	LanguageID     *uint
	ContentFile    *string
	Annotation     *map[string]string

	DOI      *string
	DOIClear *bool

	ExpertConclusionFile  *string
	ExpertConclusionClear *bool

	AffiliationsIDs *[]uint
	Tags            TagNamesByLang
}

func NewArticleUpdateEntity(name *map[string]string, coAuthorsCount *int, coAuthorsIDs *[]uint, studyFieldsIDs *[]uint, languageID *uint, contentFile *string, annotation *map[string]string, DOI *string, DOIClear *bool, expertConclusionFile *string, expertConclusionClear *bool, affiliationsIDs *[]uint) *ArticleUpdateEntity {
	return &ArticleUpdateEntity{Name: name, CoAuthorsCount: coAuthorsCount, CoAuthorsIDs: coAuthorsIDs, StudyFieldsIDs: studyFieldsIDs, LanguageID: languageID, ContentFile: contentFile, Annotation: annotation, DOI: DOI, DOIClear: DOIClear, ExpertConclusionFile: expertConclusionFile, ExpertConclusionClear: expertConclusionClear, AffiliationsIDs: affiliationsIDs}
}
