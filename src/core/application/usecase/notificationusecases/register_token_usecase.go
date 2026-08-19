package notificationusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type RegisterTokenUseCase struct {
	repository repository.NotificationTokenRepository
	gateway    gateway.NotificationGateway
}

// @inject
func NewRegisterTokenUseCase(repository repository.NotificationTokenRepository, gateway gateway.NotificationGateway) *RegisterTokenUseCase {
	return &RegisterTokenUseCase{repository: repository, gateway: gateway}
}

func (this *RegisterTokenUseCase) Execute(userID uint, token string, segment int) error {
	data := entity.NewNotificationTokenEntity(0, userID, token, enum.NotificationSegment(segment))
	isCreated, err := this.repository.CreateOrUpdate(data)

	if err != nil {
		return err
	}

	if isCreated && enum.NotificationSegment(segment) == enum.SegmentWeb {
		return this.gateway.SubscribeToTopic(enum.TopicNews, token)
	}
	return nil
}
