package doiusecases

import (
	"context"
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DoiSettingGetUseCase struct {
	repo repository.JournalDoiSettingRepository
}

// @inject
func NewDoiSettingGetUseCase(repo repository.JournalDoiSettingRepository) *DoiSettingGetUseCase {
	return &DoiSettingGetUseCase{repo: repo}
}

func (this *DoiSettingGetUseCase) Execute(ctx context.Context, journalID uint) (*entity.JournalDoiSettingEntity, error) {
	result, err := this.repo.GetByJournalID(ctx, journalID)
	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return nil, response.NewOptionalResponse(200, response.CodeNotFound, nil, "doi setting not found")
		}
		return nil, err
	}
	return result, nil
}
