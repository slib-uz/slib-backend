package mapper

import (
	"fmt"

	"slib.uz/src/core/domain/entity"
)

func JournalArticlePurchaseStatsEntityMapper(entities []*entity.JournalArticlePurchaseStatsEntity) []*entity.JournalArticlePurchaseStatsWithNameEntity {

	var result = make([]*entity.JournalArticlePurchaseStatsWithNameEntity, len(entities))
	for i, item := range entities {
		result[i] = entity.NewJournalArticlePurchaseStatsWithNameEntity(item.JournalID, item.GetJournalName(), item.Count, item.TotalAmount)
	}
	return result
}

func ArticleEntityToROIPublishEntity(article *entity.ArticleEntity, frontendURL string) *entity.ROIPublishEntity {

	var coAuthors []*entity.ROICoAuthorEntity
	for _, author := range article.CoAuthors {
		coAuthors = append(coAuthors, entity.NewROICoAuthorEntity(author.FullName, author.ScienceID))
	}

	var studyFields []*entity.ROIStudyFieldsEntity
	for _, studyField := range article.StudyFields {
		var code uint
		if studyField.Code != nil {
			code = *studyField.Code
		}
		studyFields = append(studyFields, entity.NewROIStudyFieldsEntity(studyField.Name, code))
	}

	var journal *entity.ROIJournalEntity
	if article.Journal != nil {
		journal = JournalEntityToOutput(article.Journal)
	}

	var publisher *entity.ROIPublisherEntity
	if article.Journal != nil && article.Journal.Publisher != nil {
		publisher = PublisherEntityToROIPublisherEntity(article.Journal.Publisher)
	}

	url := fmt.Sprintf("%s/article/%d", frontendURL, article.ID)

	return entity.NewROIPublishEntity(
		article.Name,
		article.PublicationDate.Format("2006-01-02"),
		article.CoAuthorsCount,
		int(article.AccessType),
		coAuthors,
		studyFields,
		article.LanguageID,
		article.Annotation,
		article.DOI,
		article.ROI,
		url,
		journal,
		publisher,
		nil,
	)
}

func ArticleEntityToROIArticleUpdateEntity(article *entity.ArticleEntity, sourceLink string) *entity.ROIArticleUpdateEntity {
	var coAuthors []*entity.ROICoAuthorEntity
	for _, author := range article.CoAuthors {
		coAuthors = append(coAuthors, entity.NewROICoAuthorEntity(author.FullName, author.ScienceID))
	}

	var studyFieldsCodes []uint
	for _, studyField := range article.StudyFields {
		if studyField.Code != nil {
			studyFieldsCodes = append(studyFieldsCodes, *studyField.Code)
		}
	}

	roiCode := ""
	if article.ROI != nil {
		roiCode = *article.ROI
	}

	return entity.NewROIArticleUpdateEntity(
		roiCode,
		article.Name,
		article.PublicationDate.Format("2006-01-02"),
		uint(article.CoAuthorsCount),
		coAuthors,
		int(article.AccessType),
		studyFieldsCodes,
		article.LanguageID,
		article.ContentFile,
		article.DOI,
		article.Annotation,
		sourceLink,
	)
}
