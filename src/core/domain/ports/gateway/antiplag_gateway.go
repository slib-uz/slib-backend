package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type AntiPlagGateway interface {
	Check(file []byte, fileName string) (*entity.AntiPlagResultEntity, error)
	GetStatus(externalID uint) int
	GetResult(externalID uint) (*entity.AntiPlagResultEntity, error)
	GetBalance() (int, error)
	GenerateCertificate(externalID uint, authorFullName, fileName string) ([]byte, error)
}
