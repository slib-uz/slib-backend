package notificationusecases

import (
	"log"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type NotificationDetailUseCase struct {
	repository repository.NotificationRepository
}

// @inject
func NewNotificationDetailUseCase(repository repository.NotificationRepository) *NotificationDetailUseCase {
	return &NotificationDetailUseCase{repository: repository}
}

func (this *NotificationDetailUseCase) Execute(userID, notificationID uint) (*entity.NotificationEntity, error) {
	n, err := this.repository.GetByID(notificationID, userID)
	if err != nil {
		return nil, err
	}

	if !n.IsRead {
		if err := this.repository.Read(n.ID, userID); err != nil {
			log.Println("Error marking notification as read:", err)
		}
	}

	return n, nil
}
