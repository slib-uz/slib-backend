package mapper

import (
	"time"

	"gorm.io/datatypes"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/persistence/models"
)

type ReviewStageOverdueRow struct {
	ReviewStageID     uint
	ApplicationID     uint
	ApplicationNumber string
	ArticleID         uint
	ArticleName       datatypes.JSON
	StageNumber       enum.Stage
	StageName         string
	Deadline          time.Time
	CreatedAt         time.Time
	OverdueDays       int
	JournalID         uint
	JournalName       datatypes.JSON
}

func ReviewStageOverdueRowsToEntities(rows []ReviewStageOverdueRow) []*entity2.ReviewStageOverdueEntity {
	entities := make([]*entity2.ReviewStageOverdueEntity, len(rows))
	for i, row := range rows {
		articleName := FromGormJson[map[string]string](row.ArticleName)
		entities[i] = &entity2.ReviewStageOverdueEntity{
			ReviewStageID:     row.ReviewStageID,
			ApplicationID:     row.ApplicationID,
			ApplicationNumber: row.ApplicationNumber,
			ArticleID:         row.ArticleID,
			Article: &entity2.ArticleEntity{
				ID:   row.ArticleID,
				Name: articleName,
			},
			StageNumber:       row.StageNumber,
			StageName:         row.StageName,
			Deadline:          row.Deadline,
			CreatedAt:         row.CreatedAt,
			OverdueDays:       row.OverdueDays,
			JournalID:         row.JournalID,
			JournalName:       FromGormJson[map[string]string](row.JournalName),
		}
	}
	return entities
}

func ReviewStageModelsToEntities(models []*models.ReviewStageModel) []*entity2.ReviewStageEntity {
	entities := make([]*entity2.ReviewStageEntity, len(models))
	for i, m := range models {
		entities[i] = ReviewStageModelToEntity(m)
	}
	return entities
}

func ReviewStageModelToEntity(s *models.ReviewStageModel) *entity2.ReviewStageEntity {
	var reviewer *entity2.UserEntity
	var application *entity2.ApplicationEntity
	var oldReviewStage *entity2.ReviewStageEntity
	var antiPlagResult *entity2.AntiPlagResultEntity
	var spellCheckResult *entity2.SpellCheckResultEntity
	var aiDetectResult *entity2.AiDetectResultEntity

	if s.Reviewer != nil {
		reviewer = UserModelToEntity(s.Reviewer)
	}
	if s.Application != nil {
		application = ApplicationModelToEntity(s.Application)
	}

	if s.Previous != nil {
		oldReviewStage = ReviewStageModelToEntity(s.Previous)
	}

	if len(s.AntiPlagResults) > 0 {
		antiPlagResult = AntiPlagResultModelToEntity(s.AntiPlagResults[0])
	}

	if len(s.SpellCheckResults) > 0 {
		spellCheckResult = SpellCheckResultModelToEntity(s.SpellCheckResults[0])
	}

	if len(s.AiDetectResults) > 0 {
		aiDetectResult = AiDetectResultModelToEntity(s.AiDetectResults[0])
	}

	return entity2.NewReviewStageEntity(
		s.ID,
		s.ApplicationID,
		application,
		s.Stage,
		s.Status,
		s.Reason,
		s.ReviewerID,
		reviewer,
		s.ReviewedAt,
		s.CreatedAt,
		s.PreviousID,
		oldReviewStage,
		s.IsOld,
		antiPlagResult,
		spellCheckResult,
		aiDetectResult,
		s.Deadline,
		s.ResubmitDeadline,
	)
}

func ReviewStateEntityToModel(s *entity2.ReviewStageEntity) models.ReviewStageModel {
	model := models.NewReviewStageModel(
		s.ID,
		s.ApplicationID,
		s.Stage,
		s.Status,
		s.PreviousID,
		s.Deadline,
		s.ResubmitDeadline,
	)

	// Set additional fields that are not in constructor
	model.Reason = s.Reason
	model.ReviewerID = s.ReviewerID
	model.ReviewedAt = s.ReviewedAt
	model.IsOld = s.IsOld

	return model
}
