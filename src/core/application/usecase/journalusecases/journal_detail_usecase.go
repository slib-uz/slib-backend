package journalusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalDetailUseCase struct {
	repository repository.JournalRepository
	viewsCache cache.JournalViewsCountCache
}

// @inject
func NewJournalDetailUseCase(repository repository.JournalRepository, viewsCache cache.JournalViewsCountCache) *JournalDetailUseCase {
	return &JournalDetailUseCase{repository: repository, viewsCache: viewsCache}
}

func (this *JournalDetailUseCase) Execute(id uint, userKey string) (*entity.JournalEntity, error) {

	journal, err := this.repository.GetByIDWithRelations(id)
	if err != nil {
		return nil, err
	}

	count, err := this.viewsCache.Add(userKey, id)
	if err == nil {
		journal.ViewsCount = journal.ViewsCount + count
	}

	return journal, nil
}
