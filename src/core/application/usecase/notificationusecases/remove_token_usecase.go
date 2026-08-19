package notificationusecases

import "slib.uz/src/core/domain/ports/repository"

type RemoveTokenUseCase struct {
	repository repository.NotificationTokenRepository
}

// @inject
func NewRemoveTokenUseCase(repository repository.NotificationTokenRepository) *RemoveTokenUseCase {
	return &RemoveTokenUseCase{repository: repository}
}

func (this *RemoveTokenUseCase) Execute(userID uint, token string) error {
	return this.repository.Delete(userID, token)
}
