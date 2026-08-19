package doiusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DoiSettingUpdateUseCase struct {
	repo repository.JournalDoiSettingRepository
}

// @inject
func NewDoiSettingUpdateUseCase(repo repository.JournalDoiSettingRepository) *DoiSettingUpdateUseCase {
	return &DoiSettingUpdateUseCase{repo: repo}
}

func (this *DoiSettingUpdateUseCase) Execute(ctx context.Context, e *entity.JournalDoiSettingEntity) error {
	return this.repo.Update(ctx, e)
}
