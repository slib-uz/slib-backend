package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type AiDetectGateway interface {
	Check(file []byte, fileName string) (*entity.AiDetectResultEntity, error)
	GetResult(externalID uint) (*entity.AiDetectResultEntity, error)
}
