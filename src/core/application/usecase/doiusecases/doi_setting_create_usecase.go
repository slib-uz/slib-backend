package doiusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DoiSettingCreateUseCase struct {
	repo repository.JournalDoiSettingRepository
}

// @inject
func NewDoiSettingCreateUseCase(repo repository.JournalDoiSettingRepository) *DoiSettingCreateUseCase {
	return &DoiSettingCreateUseCase{repo: repo}
}

func (this *DoiSettingCreateUseCase) Execute(ctx context.Context, e *entity.JournalDoiSettingEntity) error {
	return this.repo.Create(ctx, e)
}
