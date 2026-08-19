package authusecases

import (
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type TestTokenUseCase struct {
	auth       *service.UserAuthTokenService
	repository repository.UserRepository
}

// @inject
func NewTestTokenUseCase(auth *service.UserAuthTokenService, repository repository.UserRepository) *TestTokenUseCase {
	return &TestTokenUseCase{auth: auth, repository: repository}
}

func (this *TestTokenUseCase) Execute(userID uint) (*entity.AuthTokenEntity, error) {
	_, err := this.repository.GetById(userID)
	if err != nil {
		return nil, err
	}
	return this.auth.GenerateToken(userID), nil
}
