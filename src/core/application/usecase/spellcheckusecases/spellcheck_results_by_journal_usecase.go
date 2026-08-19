package spellcheckusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SpellCheckResultsByJournalUsecase struct {
	repository repository.SpellCheckResultRepository
}

// @inject
func NewSpellCheckResultsByJournalUsecase(repository repository.SpellCheckResultRepository) *SpellCheckResultsByJournalUsecase {
	return &SpellCheckResultsByJournalUsecase{
		repository: repository,
	}
}

func (this *SpellCheckResultsByJournalUsecase) Execute(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.SpellCheckResultEntity], error) {
	paging, err := this.repository.GetByJournalID(journalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*entity2.SpellCheckResultEntity, len(paging.Items))
	for i, item := range paging.Items {
		items[i] = mapper.SpellCheckResultEntityToResponseEntity(item)
	}

	return entity2.NewPagingEntity(page, pageSize, paging.Total, items), nil
}
