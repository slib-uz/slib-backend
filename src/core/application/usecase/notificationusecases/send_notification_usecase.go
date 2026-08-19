package notificationusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type SendNotificationUseCase struct {
	tokenRepo        repository.NotificationTokenRepository
	gateway          gateway.NotificationGateway
	resultRepository repository.NotificationSendResultRepository
}

// @inject
func NewSendNotificationUseCase(
	tokenRepo repository.NotificationTokenRepository,
	gateway gateway.NotificationGateway,
	resultRepository repository.NotificationSendResultRepository,
) *SendNotificationUseCase {
	return &SendNotificationUseCase{
		tokenRepo:        tokenRepo,
		gateway:          gateway,
		resultRepository: resultRepository,
	}
}

func (uc *SendNotificationUseCase) Execute(n *entity.NotificationEntity) error {
	if !n.IsBroadcast && n.UserID != nil && !n.IsEmail {
		return uc.sendPushNotification(n)
	} else if n.IsBroadcast {
		return uc.gateway.SendToTopic(n.Title["uz"], n.Body["uz"], n.ExtraData, *n.Topic)
	}
	return nil
}

func (uc *SendNotificationUseCase) sendPushNotification(n *entity.NotificationEntity) error {
	tokens, err := uc.tokenRepo.GetTokensByUserID(*n.UserID)
	if err != nil {
		return err
	}

	if tokens != nil && len(tokens) == 0 {
		return nil
	}

	result, err := uc.gateway.SendToTokens(n.ID, n.Title["uz"], n.Body["uz"], n.ExtraData, tokens)
	if err != nil {
		return err
	}

	return uc.resultRepository.Create(result)
}
