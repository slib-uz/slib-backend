package repository

import "slib.uz/src/core/domain/entity"

type SandboxUserRepository interface {
	GetByPhoneNumber(phoneNumber string) (*entity.SandboxUserEntity, error)
}
