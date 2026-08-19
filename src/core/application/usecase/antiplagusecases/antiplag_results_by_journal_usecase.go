package antiplagusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AntiPlagResultsByJournalUsecase struct {
	repository repository.AntiPlagRepository
}

// @inject
func NewAntiPlagResultsByJournalUsecase(repository repository.AntiPlagRepository) *AntiPlagResultsByJournalUsecase {
	return &AntiPlagResultsByJournalUsecase{
		repository: repository,
	}
}

func (this *AntiPlagResultsByJournalUsecase) Execute(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.AntiPlagResultEntity], error) {
	paging, err := this.repository.GetByJournalID(journalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*entity2.AntiPlagResultEntity, len(paging.Items))
	copy(items, paging.Items)

	return entity2.NewPagingEntity(page, pageSize, paging.Total, items), nil
}
