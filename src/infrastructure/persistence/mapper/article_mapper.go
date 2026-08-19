package mapper

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	. "slib.uz/src/infrastructure/persistence/models"
)

func ArticleModelListToEntityList(models []*ArticleModel) []*entity2.ArticleEntity {
	articles := make([]*entity2.ArticleEntity, len(models))
	for i, model := range models {
		articles[i] = ArticleModelToEntity(model)
	}
	return articles
}

func ArticleUpdateEntityToModel(articleID uint, oldArticle *entity2.ArticleEntity, newArticle *entity2.ArticleUpdateEntity) *ArticleModel {

	var name datatypes.JSON
	var annotation datatypes.JSON
	var file string

	coAuthorsCount := 0
	languageID := uint(0)

	if newArticle.Name != nil {
		name = ToGormJson(newArticle.Name)
	}

	if newArticle.Annotation != nil {
		annotation = ToGormJson(newArticle.Annotation)
	}

	if newArticle.LanguageID != nil {
		languageID = *newArticle.LanguageID
	}

	if newArticle.CoAuthorsCount != nil {
		coAuthorsCount = *newArticle.CoAuthorsCount
	}

	if newArticle.ContentFile != nil {
		file = *newArticle.ContentFile
	}

	return NewArticleModel(
		articleID,
		name,
		coAuthorsCount,
		nil,
		0,
		nil,
		languageID,
		annotation,
		file,
		nil,
		newArticle.DOI,
		0,
		nil,
	)
}

func ArticleCreateEntityToModel(entity *entity2.ArticleCreateEntity, accessType enum.AccessType) *ArticleModel {

	coAuthors := make([]*AuthorModel, len(entity.CoAuthorsIDs))
	for i, id := range entity.CoAuthorsIDs {
		coAuthors[i] = &AuthorModel{Model: gorm.Model{ID: id}}
	}
	studyFields := make([]*StudyFieldModel, len(entity.StudyFieldsIDs))
	for i, id := range entity.StudyFieldsIDs {
		studyFields[i] = &StudyFieldModel{Model: gorm.Model{ID: id}}
	}

	var tags []*TagModel
	for _, tagId := range entity.TagIds {
		tags = append(tags, &TagModel{ID: tagId})
	}

	return NewArticleModel(
		0,
		ToGormJson(entity.Name),
		int(entity.CoAuthorsCount),
		coAuthors,
		accessType,
		studyFields,
		entity.LanguageID,
		ToGormJson(entity.Annotation),
		entity.ContentFile,
		tags,
		entity.DOI,
		entity.JournalID,
		entity.ExpertConclusionFilePath,
	)

}

func ArticleModelToBasicEntity(model *ArticleModel) *entity2.ArticleBasicEntity {

	var coAuthors []*entity2.AuthorEntity
	var studyFields []*entity2.StudyFieldEntity
	var journal *entity2.JournalEntity

	if model.Journal != nil {
		journal = JournalModelToEntity(model.Journal)
	}

	if model.CoAuthors == nil {
		coAuthors = make([]*entity2.AuthorEntity, len(model.CoAuthors))
		for i, author := range model.CoAuthors {
			coAuthors[i] = AuthorModelToEntity(author)
		}
	}

	if model.StudyFields == nil {
		studyFields = make([]*entity2.StudyFieldEntity, len(model.StudyFields))
		for i, field := range model.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(field)
		}
	}

	ratingAvg := calculateRatingAvg(model.RatingSum, model.RatingCount)

	return entity2.NewArticleBasicEntity(
		model.ID,
		FromGormJson[map[string]string](model.Name),
		time.Time(model.PublicationDate),
		coAuthors,
		studyFields,
		model.ContentFilePath,
		model.JournalID,
		journal,
		model.DOI,
		model.ROI,
		model.RatingCount,
		model.RatingSum,
		ratingAvg,
		model.ExpertConclusionFile,
	)
}

func ArticleModelToEntity(a *ArticleModel) *entity2.ArticleEntity {

	var coAuthors []*entity2.AuthorEntity
	var studyFields []*entity2.StudyFieldEntity
	var language *entity2.LanguageEntity
	var tags entity2.TagNamesByLang
	var journal *entity2.JournalEntity
	var authorAffiliations []*entity2.ArticleAuthorAffiliationEntity

	if a.CoAuthors != nil {
		coAuthors = make([]*entity2.AuthorEntity, len(a.CoAuthors))
		for i, author := range a.CoAuthors {
			coAuthors[i] = AuthorModelToEntity(author)
		}
	}
	if a.StudyFields != nil {
		studyFields = make([]*entity2.StudyFieldEntity, len(a.StudyFields))
		for i, field := range a.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(field)
		}
	}

	if a.Language != nil {
		language = LanguageModelToEntity(a.Language)
	}

	if a.Tags != nil {
		tags = TagModelsToNamesByLang(a.Tags)
	}

	if a.Journal != nil {
		journal = JournalModelToEntity(a.Journal)
	}

	if a.ArticleAuthorAffiliations != nil {
		authorAffiliations = ArticleAuthorAffiliationListModelToEntity(a.ArticleAuthorAffiliations)
	}

	var references []*entity2.ReferenceEntity
	if a.References != nil {
		references = ReferenceModelListToEntityList(a.References)
	}

	ratingAvg := calculateRatingAvg(a.RatingSum, a.RatingCount)

	entity := entity2.NewArticleEntity(
		a.ID,
		FromGormJson[map[string]string](a.Name),
		time.Time(a.PublicationDate),
		a.CoAuthorsCount,
		coAuthors,
		a.AccessType,
		studyFields,
		a.LanguageID,
		language,
		FromGormJson[map[string]string](a.Annotation),
		a.ContentFilePath,
		tags,
		a.DOI,
		a.ROI,
		a.JournalID,
		journal,
		a.ExpertConclusionFile,
		a.ViewsCount,
		a.RatingCount,
		a.RatingSum,
		ratingAvg,
		authorAffiliations,
	)
	entity.References = references
	entity.UnconfirmedAuthors = FromGormJson[[]*entity2.UnconfirmedAuthorEntity](a.UnconfirmedAuthors)
	return entity
}

func ArticleModelToInput(a *ArticleModel) *entity2.ArticleInputEntity {
	var coAuthors []*entity2.AuthorEntity
	var studyFields []*entity2.StudyFieldEntity
	var language *entity2.LanguageEntity
	var tags entity2.TagNamesByLang
	var journal *entity2.JournalEntity
	var authorAffiliations []*entity2.ArticleAuthorAffiliationEntity

	if a.CoAuthors != nil {
		coAuthors = make([]*entity2.AuthorEntity, len(a.CoAuthors))
		for i, author := range a.CoAuthors {
			coAuthors[i] = AuthorModelToEntity(author)
		}
	}
	if a.ArticleAuthorAffiliations != nil {
		authorAffiliations = ArticleAuthorAffiliationListModelToEntity(a.ArticleAuthorAffiliations)
	}
	if a.StudyFields != nil {
		studyFields = make([]*entity2.StudyFieldEntity, len(a.StudyFields))
		for i, field := range a.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(field)
		}
	}

	if a.Language != nil {
		language = LanguageModelToEntity(a.Language)
	}

	if a.Tags != nil {
		tags = TagModelsToNamesByLang(a.Tags)
	}

	if a.Journal != nil {
		journal = JournalModelToEntity(a.Journal)
	}

	var references []*entity2.ReferenceEntity
	if a.References != nil {
		references = ReferenceModelListToEntityList(a.References)
	}

	ratingAvg := calculateRatingAvg(a.RatingSum, a.RatingCount)

	entity := entity2.NewArticleInputEntity(
		a.ID,
		FromGormJson[map[string]string](a.Name),
		time.Time(a.PublicationDate),
		a.CoAuthorsCount,
		coAuthors,
		int(a.AccessType),
		studyFields,
		a.LanguageID,
		language,
		FromGormJson[map[string]string](a.Annotation),
		a.ContentFilePath,
		tags,
		a.DOI,
		a.ROI,
		a.JournalID,
		journal,
		a.ExpertConclusionFile,
		a.ViewsCount,
		a.RatingCount,
		a.RatingSum,
		ratingAvg,
	)
	entity.References = references
	entity.CoAuthorsWithAffiliation = authorAffiliations
	entity.UnconfirmedAuthors = FromGormJson[[]*entity2.UnconfirmedAuthorEntity](a.UnconfirmedAuthors)
	return entity
}

func ArticleModelToOnlyNameEntity(a *ArticleModel) *entity2.ArticleOnlyNameEntity {
	return entity2.NewArticleOnlyNameEntity(a.ID, FromGormJson[map[string]string](a.Name))
}

func ArticleModelListToOnlyNameEntityList(models []*ArticleModel) []*entity2.ArticleOnlyNameEntity {
	entities := make([]*entity2.ArticleOnlyNameEntity, len(models))
	for i, model := range models {
		entities[i] = ArticleModelToOnlyNameEntity(model)
	}
	return entities
}

func ArticlePublishEntityToArticleModel(e *entity2.ArticlePublishEntity, journalID uint) *ArticleModel {
	accessType := e.AccessType
	if accessType == 0 {
		accessType = enum.PrivateAccessType
	}

	article := NewArticleModel(
		0,
		ToGormJson(e.Name),
		int(e.CoAuthorsCount),
		nil,
		accessType,
		nil,
		e.LanguageID,
		ToGormJson(e.Annotation),
		e.ContentFile,
		nil,
		e.DOI,
		journalID,
		e.ExpertConclusionFile,
	)
	article.UnconfirmedAuthors = ToGormJson(e.UnconfirmedAuthors)
	return article
}
