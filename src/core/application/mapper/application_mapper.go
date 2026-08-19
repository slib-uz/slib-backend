package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func ApplicationEntityListToBasicEntityList(applications []*entity.ApplicationEntity) []*entity.ApplicationBasicEntity {
	entities := make([]*entity.ApplicationBasicEntity, len(applications))
	for i, application := range applications {
		entities[i] = ApplicationEntityToBasicEntity(application)
	}
	return entities
}

func ApplicationEntityToBasicEntity(application *entity.ApplicationEntity) *entity.ApplicationBasicEntity {
	var journal *entity.JournalEntity
	var currentStage *entity.ReviewStageResponseEntity

	if application.Journal != nil {
		journal = application.Journal
	}

	if len(application.ReviewStages) > 0 {
		currentStage = ReviewStageEntityToResponseEntity(application.ReviewStages[0])
	}

	return entity.NewApplicationBasicEntity(
		application.ID,
		application.Number,
		application.ArticleID,
		application.Article,
		application.JournalID,
		journal,
		application.UserID,
		application.User,
		currentStage,
	)
}

func ApplicationEntityToResponse(application *entity.ApplicationEntity) *entity.ApplicationResponseEntity {
	var article *entity.ArticleInputEntity
	var journal *entity.JournalEntity
	var reviewStages []*entity.ReviewStageResponseEntity

	if application.Article != nil {
		article = ArticleEntityToInput(application.Article)
		slimApplicationCoAuthors(article)
	}
	if application.Journal != nil {
		journal = application.Journal
	}

	if application.ReviewStages != nil {
		reviewStages = make([]*entity.ReviewStageResponseEntity, len(application.ReviewStages))
		for i, stage := range application.ReviewStages {
			reviewStages[i] = ReviewStageEntityToResponseEntity(stage)
		}
	}

	return entity.NewApplicationResponseEntity(
		application.ID,
		application.Number,
		application.ArticleID,
		article,
		application.JournalID,
		journal,
		application.UserID,
		application.User,
		reviewStages,
		application.CreatedAt,
		application.UpdatedAt,
	)
}

// slimApplicationCoAuthors keeps only id, full_name, science_id, orcid_id on co_authors
// (drops nested user). Affiliation authors stay as loaded for co_authors_with_affiliation.
func slimApplicationCoAuthors(article *entity.ArticleInputEntity) {
	if article == nil {
		return
	}
	for _, author := range article.CoAuthors {
		if author == nil {
			continue
		}
		author.User = nil
	}
}
