package service

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ClientAuthService struct {
	clientRepo repository.ClientRepository
}

// @inject
func NewClientAuthService(clientRepo repository.ClientRepository) *ClientAuthService {
	return &ClientAuthService{
		clientRepo: clientRepo,
	}
}

func (s *ClientAuthService) Authenticate(clientID, clientSecret string) (*entity.ClientEntity, error) {
	client, err := s.clientRepo.GetByClientID(clientID)
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, errors.New("client not found")
	}

	if !client.IsActive {
		return nil, errors.New("client is inactive")
	}

	if clientSecret != client.ClientSecret {
		return nil, errors.New("invalid client secret")
	}

	return client, nil
}
