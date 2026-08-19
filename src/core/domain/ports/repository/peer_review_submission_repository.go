package repository

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type PeerReviewSubmissionRepository interface {
	BulkCreate([]*entity.PeerReviewSubmissionEntity) error
	UpdateStatusAndDeadline(externalID uint, status enum.PeerReviewStatus, deadline time.Time, oldDeadline *time.Time, version int64) error
	GetByExternalID(externalID uint) (*entity.PeerReviewSubmissionEntity, error)
	GetByApplicationID(applicationID uint) ([]*entity.PeerReviewSubmissionEntity, error)
	GetByIDWithApplication(id uint) (*entity.PeerReviewSubmissionEntity, error)
	IsExistsByExternalIdempotencyKey(idempotencyKey string) (bool, error)
}
