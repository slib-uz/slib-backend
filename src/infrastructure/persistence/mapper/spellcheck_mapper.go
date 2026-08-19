package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func SpellCheckResultEntityToModel(data *entity.SpellCheckResultEntity) *models.SpellCheckResultModel {

	return models.NewSpellCheckResultModel(
		data.ReviewStageID,
		data.ApplicationID,
		data.JournalID,
		data.File,
		data.ResultFile,
		data.Status,
		data.SubmitterID,
	)
}

func SpellCheckResultModelToEntity(data *models.SpellCheckResultModel) *entity.SpellCheckResultEntity {

	var application *entity.ApplicationBasicEntity
	var user *entity.UserBasicEntity
	var journal *entity.JournalEntity
	var reviewStage *entity.ReviewStageResponseEntity

	if data.Application != nil {
		application = ApplicationModelToBasicEntity(data.Application)
	}

	if data.Submitter != nil {
		user = UserModelToBasic(data.Submitter)
	}

	if data.Journal != nil {
		journal = JournalModelToEntity(data.Journal)
	}

	if data.ReviewStage != nil {
		reviewStage = ReviewStageModelToResponseEntity(data.ReviewStage)
	}

	return entity.NewSpellCheckResultEntity(
		data.ID,
		data.ReviewStageID,
		reviewStage,
		data.ApplicationID,
		application,
		data.JournalID,
		journal,
		data.File,
		data.ResultFile,
		data.Status,
		data.SubmitterID,
		user,
		data.ResultTime,
	)
}

func ApplicationModelToBasicEntity(model *models.ArticleApplicationModel) *entity.ApplicationBasicEntity {
	appEntity := ApplicationModelToEntity(model)
	return ApplicationEntityToBasicEntity(appEntity)
}

func ApplicationEntityToBasicEntity(it *entity.ApplicationEntity) *entity.ApplicationBasicEntity {
	var journal *entity.JournalEntity
	var currentStage *entity.ReviewStageResponseEntity

	if it.Journal != nil {
		journal = it.Journal
	}

	if len(it.ReviewStages) > 0 {
		currentStage = ReviewStageEntityToResponseEntity(it.ReviewStages[0])
	}

	return entity.NewApplicationBasicEntity(
		it.ID,
		it.Number,
		it.ArticleID,
		it.Article,
		it.JournalID,
		journal,
		it.UserID,
		it.User,
		currentStage,
	)
}

func UserModelToBasic(model *models.UserModel) *entity.UserBasicEntity {
	userEntity := UserModelToEntity(model)
	return UserEntityToBasic(userEntity)
}

func UserEntityToBasic(user *entity.UserEntity) *entity.UserBasicEntity {
	return entity.NewUserBasicEntity(
		user.ID,
		user.ScienceID,
		user.FullName,
		user.PhoneNumber,
		user.IsAdmin,
		user.Roles,
	)
}

func ReviewStageModelToResponseEntity(model *models.ReviewStageModel) *entity.ReviewStageResponseEntity {
	reviewStageEntity := ReviewStageModelToEntity(model)
	return ReviewStageEntityToResponseEntity(reviewStageEntity)
}

func ReviewStageEntityToResponseEntity(s *entity.ReviewStageEntity) *entity.ReviewStageResponseEntity {
	var application *entity.ApplicationBasicEntity
	var reviewer *entity.UserBasicEntity
	var previous *entity.ReviewStageResponseEntity
	var antiPlagResult *entity.AntiPlagResultEntity
	var spellCheckResult *entity.SpellCheckResultEntity
	var aiDetectResult *entity.AiDetectResultEntity

	if s.Application != nil {
		application = ApplicationEntityToBasicEntity(s.Application)
	}
	if s.Reviewer != nil {
		reviewer = UserEntityToBasic(s.Reviewer)
	}
	if s.Previous != nil {
		previous = ReviewStageEntityToResponseEntity(s.Previous)
	}

	if s.AntiPlagResult != nil {
		antiPlagResult = s.AntiPlagResult
	}

	if s.SpellCheckResult != nil {
		spellCheckResult = s.SpellCheckResult
	}

	if s.AIDetectResult != nil {
		aiDetectResult = s.AIDetectResult
	}

	return entity.NewReviewStageResponseEntity(
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
		previous,
		antiPlagResult,
		spellCheckResult,
		aiDetectResult,
		s.IsOld,
		s.Deadline,
	)
}
