package aidetectusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AiDetectResultsByJournalUsecase struct {
	repository repository.AiDetectRepository
}

// @inject
func NewAiDetectResultsByJournalUsecase(repository repository.AiDetectRepository) *AiDetectResultsByJournalUsecase {
	return &AiDetectResultsByJournalUsecase{repository: repository}
}

func (this *AiDetectResultsByJournalUsecase) Execute(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.AiDetectResultEntity], error) {
	paging, err := this.repository.GetByJournalID(journalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*entity2.AiDetectResultEntity, len(paging.Items))
	copy(items, paging.Items)

	return entity2.NewPagingEntity(page, pageSize, paging.Total, items), nil
}
