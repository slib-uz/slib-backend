package repository

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/config"
)

type DeveloperUserRepositoryImpl struct {
	config *config.Config
}

// @inject
func NewDeveloperUserRepository(config *config.Config) repository.DeveloperUserRepository {
	return &DeveloperUserRepositoryImpl{config: config}
}

func (this *DeveloperUserRepositoryImpl) GetByUsername(username string) (*entity.DeveloperUserEntity, error) {
	if username != this.config.DeveloperUsername {
		return nil, errors.New("user not found")
	}

	return &entity.DeveloperUserEntity{
		ID:       1,
		Username: this.config.DeveloperUsername,
		Password: this.config.DeveloperPasswordHash,
	}, nil
}

func (this *DeveloperUserRepositoryImpl) GetByID(id uint) (*entity.DeveloperUserEntity, error) {
	return &entity.DeveloperUserEntity{
		ID:       1,
		Username: this.config.DeveloperUsername,
		Password: this.config.DeveloperPasswordHash,
	}, nil
}
