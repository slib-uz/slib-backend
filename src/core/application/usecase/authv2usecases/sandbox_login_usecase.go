package authv2usecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/unitofwork"
)

type SandboxLoginUseCase struct {
	sandboxUserRepository repository.SandboxUserRepository
	userRepository        repository.UserRepository
	profileRepository     repository.UserProfileRepository
	atomic                unitofwork.Atomic
	authTokenService      *service.UserAuthTokenService
}

// @inject
func NewSandboxLoginUseCase(sandboxUserRepository repository.SandboxUserRepository, userRepository repository.UserRepository, profileRepository repository.UserProfileRepository, atomic unitofwork.Atomic, authTokenService *service.UserAuthTokenService) *SandboxLoginUseCase {
	return &SandboxLoginUseCase{sandboxUserRepository: sandboxUserRepository, userRepository: userRepository, profileRepository: profileRepository, atomic: atomic, authTokenService: authTokenService}
}

func (this *SandboxLoginUseCase) Execute(phone string, otp string) (*entity.AuthTokenEntity, error) {
	sandboxUser, err := this.sandboxUserRepository.GetByPhoneNumber(phone)
	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return nil, response.NotFound
		}
		return nil, err
	}

	if sandboxUser.Otp != otp {
		return nil, response.NewOptionalResponse(200, response.CodeInvalidOrReusedCode, nil, "Invalid OTP")
	}

	user, err := this.userRepository.GetByPhoneNumber(phone)
	if err != nil && !errors.Is(err, response.NotFoundError) {
		return nil, err
	}
	if user != nil {
		return this.authTokenService.GenerateToken(user.ID), nil
	}

	var userID uint
	if err := this.atomic.Transaction(func(tx unitofwork.Tx) error {
		userID, err = this.userRepository.CreateByPhoneNumber(tx, &entity.UserEntity{PhoneNumber: phone, FullName: sandboxUser.FullName, ScienceID: sandboxUser.ScienceID})
		if err != nil {
			return err
		}
		return this.profileRepository.Create(tx, userID)
	}); err != nil {
		return nil, err
	}

	return this.authTokenService.GenerateToken(userID), nil

}
