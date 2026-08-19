package repository

import (
	"context"

	"github.com/google/uuid"
	"slib.uz/src/core/domain/entity"
)

type TelegramAuthSessionRepository interface {
	Create(ctx context.Context, session *entity.TelegramAuthSessionEntity) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*entity.TelegramAuthSessionEntity, error)
	Update(ctx context.Context, session *entity.TelegramAuthSessionEntity) error
	GetWaitingSessionByChatID(ctx context.Context, chatID int64) (*entity.TelegramAuthSessionEntity, error)
}
