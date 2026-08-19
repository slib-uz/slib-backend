package repository

import (
	"slib.uz/src/core/domain/entity"
)

type UserRoleRepository interface {
	GetByUserID(userID uint) ([]*entity.UserRoleEntity, error)
	Create(userRole *entity.UserRoleEntity) error
	Exists(userID uint, publisherID, journalID, institutionID *uint) (bool, error)
	GetByPublisherID(publisherID uint) ([]*entity.UserRoleEntity, error)
	GetByInstitutionID(institutionID uint) ([]*entity.UserRoleEntity, error)
	Delete(id uint) error
	GetByJournalID(journalID uint) ([]*entity.UserRoleEntity, error)
	FindByJournalID(journalID uint) ([]*entity.UserRoleEntity, error)
	GetJournalMemberIds(journalID uint) ([]uint, error)
	GetByJournalIDOrderByLastOnline(journalID uint) ([]*entity.JournalMemberLastOnlineEntity, error)
}
