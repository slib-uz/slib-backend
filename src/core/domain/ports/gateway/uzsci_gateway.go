package gateway

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

type UzsciGateway interface {
	GetAllForms(ratingPeriodID uint) ([]*entity.UzSciFormEntity, error)
	GetPublicRatingPeriods(isActive *bool) ([]*entity.UzSciRatingPeriodEntity, error)
	GetJournalByISSN(issnPaper, issnOnline string) (*response.UzSciJournalData, error)
	CreateApplication(periodID uint, journalID uint, answers []entity.UzSciApplicationAnswerEntity) error
}
