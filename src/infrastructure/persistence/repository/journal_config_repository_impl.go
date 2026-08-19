package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JournalConfigRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalConfigRepositoryImpl(baseRepository *BaseRepository) repository.JournalConfigRepository {
	return &JournalConfigRepositoryImpl{BaseRepository: baseRepository}
}

func (this *JournalConfigRepositoryImpl) ExistsByDomain(domainVariants []string) (bool, error) {
	var count int64
	err := this.db().Model(&models.JournalConfigModel{}).Where("website_url IN ?", domainVariants).Count(&count).Error
	return count > 0, infraError.Wrap(err)
}

func (this *JournalConfigRepositoryImpl) GetByWebsiteURL(s string) (*entity.JournalConfigEntity, error) {
	var _model models.JournalConfigModel
	if err := this.db().Where("website_url = ?", s).Preload("Journal").Preload("Creator").First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.JournalConfigModelToEntity(&_model), nil
}

func (this *JournalConfigRepositoryImpl) CreateOrUpdate(it *entity.JournalConfigEntity) error {
	var _model = mapper.JournalConfigEntityToModel(it)

	return infraError.Wrap(this.db().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "journal_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"creator_id", "website_url", "conf", "is_active", "updated_at"}),
	}).Create(_model).Error)

}

func (this *JournalConfigRepositoryImpl) Update(it *entity.JournalConfigEntity) error {
	//TODO implement me
	panic("implement me")
}

func (this *JournalConfigRepositoryImpl) GetByJournalID(id uint) (*entity.JournalConfigEntity, error) {
	var _model models.JournalConfigModel
	if err := this.db().Where("journal_id = ?", id).Preload("Journal").Preload("Creator").First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.JournalConfigModelToEntity(&_model), nil
}

func (this *JournalConfigRepositoryImpl) List(creatorID, journalID uint, isActive *bool, page, pageSize int) (*entity.PagingEntity[entity.JournalConfigEntity], error) {
	var _models []models.JournalConfigModel
	query := this.db().Preload("Journal", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "pin", "full_name")
	})

	if creatorID != 0 {
		query = query.Where("creator_id = ?", creatorID)
	}
	if journalID != 0 {
		query = query.Where("journal_id = ?", journalID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	var total int64
	if err := query.Model(&models.JournalConfigModel{}).Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	entities := make([]*entity.JournalConfigEntity, len(_models))
	for i, m := range _models {
		entities[i] = mapper.JournalConfigModelToEntity(&m)
	}

	return entity.NewPagingEntity(page, pageSize, total, entities), nil
}
