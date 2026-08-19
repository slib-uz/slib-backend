package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
)

func ReviewStageEntityToResponseEntity(s *entity2.ReviewStageEntity) *entity2.ReviewStageResponseEntity {

	var application *entity2.ApplicationBasicEntity
	var reviewer *entity2.UserBasicEntity
	var previous *entity2.ReviewStageResponseEntity
	var spellCheckResult *entity2.SpellCheckResultEntity
	var aiDetectResult *entity2.AiDetectResultEntity

	if s.Application != nil {
		application = ApplicationEntityToBasicEntity(s.Application)
	}
	if s.Reviewer != nil {
		reviewer = UserEntityToBasic(s.Reviewer)
	}
	if s.Previous != nil {
		previous = ReviewStageEntityToResponseEntity(s.Previous)
	}

	if s.SpellCheckResult != nil {
		spellCheckResult = SpellCheckResultEntityToResponseEntity(s.SpellCheckResult)
	}

	if s.AIDetectResult != nil {
		aiDetectResult = s.AIDetectResult
	}

	return entity2.NewReviewStageResponseEntity(
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
		s.AntiPlagResult,
		spellCheckResult,
		aiDetectResult,
		s.IsOld,
		s.Deadline,
	)
}
