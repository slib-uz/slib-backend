package authv2usecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
)

type CheckPhoneNumberUseCase struct {
	repository repository.UserRepository
}

// @inject
func NewCheckPhoneNumberUseCase(repository repository.UserRepository) *CheckPhoneNumberUseCase {
	return &CheckPhoneNumberUseCase{repository: repository}
}

func (this *CheckPhoneNumberUseCase) Execute(phoneNumber string) error {
	if !utils.IsValidPhoneNumber(phoneNumber) {
		return response.InvalidPhoneNumberError
	}

	_, err := this.repository.GetByPhoneNumber(phoneNumber)

	if errors.Is(err, response.NotFoundError) {
		return response.NewOptionalResponse(200, response.CodeNotFound, nil, "User not found by this phone number")
	}
	if err != nil {
		return err
	}

	return response.NewOptionalResponse(200, response.CodeAlreadyExists, nil, "User already exists with this phone number")
}
