package permissionusecases

import "slib.uz/src/core/domain/ports/repository"

type ApplicationOwnerPermissionUseCase struct {
	repository repository.ApplicationRepository
}

// @inject
func NewApplicationOwnerPermissionUseCase(repository repository.ApplicationRepository) *ApplicationOwnerPermissionUseCase {
	return &ApplicationOwnerPermissionUseCase{repository: repository}
}

func (this *ApplicationOwnerPermissionUseCase) Execute(userID, applicationID uint) (bool, error) {
	application, err := this.repository.FindByID(applicationID)
	if err != nil {
		return false, err
	}

	return application.UserID == userID, nil
}
