package journal_applications_usecases

import (
	"fmt"

	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationsListUsecase struct {
	repository repository.JournalApplicationRepository
}

// @inject
func NewApplicationsListUsecase(repository repository.JournalApplicationRepository) *ApplicationsListUsecase {
	return &ApplicationsListUsecase{repository: repository}
}

func (this *ApplicationsListUsecase) Execute(user *entity2.UserBasicEntity, publisherID uint, page, size int, status int) (*entity2.PagingEntity[entity2.JournalApplicationEntity], error) {

	if !this.checkAccess(user, publisherID) {
		return nil, response.PermissionDeniedError
	}

	paging, err := this.repository.GetListByPaging(publisherID, page, size, enum.StatusFromInt(status))

	if err != nil {
		return nil, err
	}

	return paging, nil

}

func (this *ApplicationsListUsecase) getStatus(status *int) *enum.Status {
	if status == nil {
		return nil
	}
	s := enum.Status(*status)
	return &s
}

func (this *ApplicationsListUsecase) checkAccess(user *entity2.UserBasicEntity, publisherID uint) bool {
	fmt.Println("Is admin: ", user.IsAdmin)

	if user.IsAdmin {
		return true
	}

	for _, role := range user.Roles {
		if role.PublisherID != nil && *role.PublisherID == publisherID {
			return true
		}
	}
	return false
}
