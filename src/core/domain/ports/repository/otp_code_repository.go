package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
)

type OTPCodeRepository interface {
	Create(ctx context.Context, otp *entity.OTPCodeEntity) (uint, error)
	GetByCodeAndSessionID(ctx context.Context, sessionID string, code string) (*entity.OTPCodeEntity, error)
	MarkAsUsed(ctx context.Context, id uint) error
}
