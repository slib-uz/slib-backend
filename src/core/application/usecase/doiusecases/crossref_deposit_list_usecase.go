package doiusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type CrossRefDepositListUseCase struct {
	doiDepositRepo repository.DoiDepositRepository
}

// @inject
func NewCrossRefDepositListUseCase(doiDepositRepo repository.DoiDepositRepository) *CrossRefDepositListUseCase {
	return &CrossRefDepositListUseCase{doiDepositRepo: doiDepositRepo}
}

func (this *CrossRefDepositListUseCase) Execute(ctx context.Context, user *entity.UserBasicEntity, journalID uint, page, pageSize int) (*entity.PagingEntity[entity.DoiDepositEntity], error) {
	var journalIDs []uint

	if !user.IsAdmin {
		for _, role := range user.Roles {
			if role.JournalID != nil {
				journalIDs = append(journalIDs, *role.JournalID)
			}
		}
	}

	if journalID > 0 {
		journalIDs = []uint{journalID}
	} else if !user.IsAdmin && len(journalIDs) == 0 {
		journalIDs = []uint{0}
	}

	return this.doiDepositRepo.List(ctx, journalIDs, page, pageSize)
}
