package mapper

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

func ArticlePublishEntityToCreate(dto *entity2.ArticlePublishEntity, journalID uint, tagIds []uint) *entity2.ArticleCreateEntity {

	return entity2.NewArticleCreateEntity(
		dto.Name,
		dto.CoAuthorsCount,
		dto.CoAuthorsIDs,
		dto.StudyFieldsIDs,
		dto.LanguageID,
		dto.ContentFile,
		tagIds,
		dto.DOI,
		journalID,
		dto.ExpertConclusionFile,
		dto.Annotation,
		dto.References,
	)
}

func ArticleEntityToInput(ent *entity2.ArticleEntity) *entity2.ArticleInputEntity {
	var journal *entity2.JournalEntity

	if ent.Journal != nil {
		journal = ent.Journal
	}

	entity := entity2.NewArticleInputEntity(
		ent.ID,
		ent.Name,
		ent.PublicationDate,
		ent.CoAuthorsCount,
		ent.CoAuthors,
		int(ent.AccessType),
		ent.StudyFields,
		ent.LanguageID,
		ent.Language,
		ent.Annotation,
		ent.ContentFile,
		ent.Tags,
		ent.DOI,
		ent.ROI,
		ent.JournalID,
		journal,
		ent.ExpertConclusionFile,
		ent.ViewsCount,
		ent.RatingCount,
		ent.RatingSum,
		ent.RatingAvg,
	)

	entity.References = ent.References
	entity.CoAuthorsWithAffiliation = ent.ArticleAuthorAffiliation
	entity.UnconfirmedAuthors = ent.UnconfirmedAuthors

	return entity
}

func ArticleEntityListToInput(entities []*entity2.ArticleEntity) []*entity2.ArticleInputEntity {
	inputs := make([]*entity2.ArticleInputEntity, len(entities))
	for i, entity := range entities {
		inputs[i] = ArticleEntityToInput(entity)
	}
	return inputs
}

func ArticleEntityToPublish(dto *entity2.ArticlePublishEntity) *entity2.ArticleEntity {
	return entity2.NewArticleEntity(
		0,
		dto.Name,
		time.Now(),
		int(dto.CoAuthorsCount),
		[]*entity2.AuthorEntity{},
		enum.AccessType(0),
		[]*entity2.StudyFieldEntity{},
		dto.LanguageID,
		nil,
		dto.Annotation,
		dto.ContentFile,
		nil,
		dto.DOI,
		nil,
		0,
		nil,
		dto.ExpertConclusionFile,
		0,
		0,
		0,
		0,
		nil,
	)
}

func ArticleEntityToBasicEntity(e *entity2.ArticleEntity) *entity2.ArticleBasicEntity {
	return entity2.NewArticleBasicEntity(
		e.ID,
		e.Name,
		e.PublicationDate,
		e.CoAuthors,
		e.StudyFields,
		e.ContentFile,
		e.JournalID,
		e.Journal,
		e.DOI,
		e.ROI,
		e.RatingCount,
		e.RatingSum,
		e.RatingAvg,
		e.ExpertConclusionFile,
	)
}
