package notificationusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserNotificationsListUseCase struct {
	repository repository.NotificationRepository
}

// @inject
func NewUserNotificationsListUseCase(repository repository.NotificationRepository) *UserNotificationsListUseCase {
	return &UserNotificationsListUseCase{repository: repository}
}

func (this *UserNotificationsListUseCase) Execute(userID uint, unread bool, page int, pageSize int) (*entity2.PagingEntity[entity2.NotificationEntity], error) {

	var paging *entity2.PagingEntity[entity2.NotificationEntity]
	var err error

	if unread {
		paging, err = this.repository.GetUnreadNotificationsByUserID(userID, page, pageSize)
	} else {
		paging, err = this.repository.GetByUserIDAndBroadcast(userID, page, pageSize)
	}

	if err != nil {
		return nil, err
	}

	return paging, nil
}
