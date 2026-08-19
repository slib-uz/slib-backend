package notificationusecases

import "slib.uz/src/core/domain/ports/repository"

type UnreadCountUseCase struct {
	repository repository.NotificationRepository
}

// @inject
func NewUnreadCountUseCase(repository repository.NotificationRepository) *UnreadCountUseCase {
	return &UnreadCountUseCase{repository: repository}
}

func (this *UnreadCountUseCase) Execute(userID uint) (int64, error) {
	return this.repository.GetUnreadCount(userID)
}
