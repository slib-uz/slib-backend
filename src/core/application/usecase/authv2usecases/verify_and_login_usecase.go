package authv2usecases

import (
	"context"
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/session"
)

type VerifyAndLoginUseCase struct {
	repository        repository.UserRepository
	profileRepository repository.UserProfileRepository
	service           *service.OTPService
	authTokenService  *service.UserAuthTokenService
	atomic            session.Atomic
	throttle          *service.AuthThrottle
}

// @inject
func NewVerifyAndLoginUseCase(repository repository.UserRepository, profileRepository repository.UserProfileRepository, service *service.OTPService, authTokenService *service.UserAuthTokenService, atomic session.Atomic, throttle *service.AuthThrottle) *VerifyAndLoginUseCase {
	return &VerifyAndLoginUseCase{repository: repository, profileRepository: profileRepository, service: service, authTokenService: authTokenService, atomic: atomic, throttle: throttle}
}

func (this *VerifyAndLoginUseCase) Execute(ctx context.Context, sessionID string, code string, scienceID, fullName string) (*entity.AuthTokenEntity, enum.AuthScope, error) {
	if this.throttle.CheckAndHitOTPVerify(ctx, sessionID) {
		return nil, "", response.TooManyAttemptsError
	}

	otp, err := this.service.Check(ctx, code, sessionID)
	if err != nil {
		return nil, "", err
	}

	this.throttle.ResetOTPVerify(ctx, sessionID)

	user, err := this.repository.GetByPhoneNumber(otp.Phone)
	if user != nil {
		return this.authTokenService.GenerateToken(user.ID), enum.AuthScopeLogin, nil
	}

	if errors.Is(err, response.NotFoundError) {
		userID, err := this.register(otp.Phone, scienceID, fullName)
		if err != nil {
			return nil, "", err
		}
		return this.authTokenService.GenerateToken(userID), enum.AuthScopeRegister, nil
	}

	return nil, "", err
}

func (this *VerifyAndLoginUseCase) register(phone, scienceID, fullName string) (uint, error) {
	var userID uint

	_user, err := this.repository.GetByScienceId(scienceID)
	if err != nil && !errors.Is(err, response.NotFoundError) {
		return 0, err
	}

	if _user != nil {
		return 0, response.NewOptionalResponse(200, response.CodeAlreadyExists, map[string]any{
			"phone_number": _user.PhoneNumber,
			"science_id":   _user.ScienceID,
		}, "This Science ID is already registered with another user")
	}

	err = this.atomic.Transaction(func(tx session.Tx) error {
		user := &entity.UserEntity{PhoneNumber: phone, ScienceID: scienceID, FullName: fullName}
		id, err := this.repository.CreateByPhoneNumber(tx, user)
		if err != nil {
			return err
		}
		userID = id

		if err := this.profileRepository.Create(tx, userID); err != nil {
			return err
		}

		return nil
	})

	return userID, err
}
