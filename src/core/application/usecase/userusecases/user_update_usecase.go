package userusecases

import (
	"fmt"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type UserUpdateUseCase struct {
	repository repository.UserRepository
	storage    storage.FileStorage
}

// @inject
func NewUserUpdateUseCase(repository repository.UserRepository, storage storage.FileStorage) *UserUpdateUseCase {
	return &UserUpdateUseCase{
		repository: repository,
		storage:    storage,
	}
}

func (this *UserUpdateUseCase) Execute(userID uint, updateDTO *entity.UserUpdateEntity) (*entity.UserUpdateEntity, error) {
	user, err := this.repository.GetById(userID)
	if err != nil {
		return nil, err
	}

	updateUser := *user

	if updateDTO.Photo != nil {
		updateUser.Photo = *updateDTO.Photo
	}

	if updateDTO.Email != nil {
		updateUser.Email = *updateDTO.Email
	}
	if updateDTO.AcademicDegree != nil {
		degreeCode := string(*updateDTO.AcademicDegree)
		if degreeCode != "" {
			if _, exists := enum.GetDegreeByCode(degreeCode); !exists {
				return nil, fmt.Errorf("invalid academic degree code: %s", degreeCode)
			}
		}
		updateUser.AcademicDegree = *updateDTO.AcademicDegree
	}
	if updateDTO.AcademicTitle != nil {
		updateUser.AcademicTitle = updateDTO.AcademicTitle
	}
	if updateDTO.ORCIDID != nil {
		updateUser.ORCIDID = updateDTO.ORCIDID
	}

	updatedUser, err := this.repository.Update(userID, &updateUser)
	if err != nil {
		return nil, err
	}

	updateUserEntity := entity.NewUserUpdateEntity(
		&updatedUser.Photo,
		&updatedUser.Email,
		&updatedUser.AcademicDegree,
		updatedUser.AcademicTitle,
		updatedUser.ORCIDID,
	)

	return updateUserEntity, nil
}
