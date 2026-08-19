package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type DoiDepositRepository interface {
	Create(ctx context.Context, entity *entity.DoiDepositEntity) error
	GetByBatchID(ctx context.Context, batchID string) (*entity.DoiDepositEntity, error)
	UpdateByBatchID(ctx context.Context, batchID string, status enum.DoiDepositStatus, message string, submissionID string, requestBody string, responseBody string) error
	List(ctx context.Context, journalIDs []uint, page, pageSize int) (*entity.PagingEntity[entity.DoiDepositEntity], error)
}
