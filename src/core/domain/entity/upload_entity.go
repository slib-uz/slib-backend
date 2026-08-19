package entity

import (
	"github.com/google/uuid"
	"slib.uz/src/core/domain/entity/enum"
)

type UploadEntity struct {
	ID         uuid.UUID
	Path       string
	AccessType enum.AccessType
}

func NewUploadEntity(ID uuid.UUID, path string, accessType enum.AccessType) *UploadEntity {
	return &UploadEntity{ID: ID, Path: path, AccessType: accessType}
}
