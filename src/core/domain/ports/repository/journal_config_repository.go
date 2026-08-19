package repository

import "slib.uz/src/core/domain/entity"

type JournalConfigRepository interface {
	CreateOrUpdate(*entity.JournalConfigEntity) error
	Update(*entity.JournalConfigEntity) error
	GetByJournalID(uint) (*entity.JournalConfigEntity, error)
	GetByWebsiteURL(string) (*entity.JournalConfigEntity, error)
	List(creatorID, journalID uint, isActive *bool, page, pageSize int) (*entity.PagingEntity[entity.JournalConfigEntity], error)
	ExistsByDomain(domainVariants []string) (bool, error)
}
