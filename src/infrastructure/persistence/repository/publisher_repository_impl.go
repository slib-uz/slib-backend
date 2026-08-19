package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type PublisherRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewPublisherRepository(repository *BaseRepository) repository.PublisherRepository {
	return &PublisherRepositoryImpl{BaseRepository: repository}
}

func (this *PublisherRepositoryImpl) GetListByPage(page, size int, name, tin string, institutionID uint, unassigned bool) (*entity2.PagingEntity[entity2.PublisherEntity], error) {

	var publishModel []*models.PublisherModel

	query := this.db().Model(&models.PublisherModel{})

	if name != "" {
		query = query.Where("name ILIKE ? OR short_name ILIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if tin != "" {
		query = query.Where("tin ILIKE ?", "%"+tin+"%")
	}

	if unassigned {
		query = query.Where("institution_id IS NULL")
	}

	if institutionID > 0 {
		query = query.Where("institution_id = ?", institutionID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Limit(size).Offset((page - 1) * size).Order("created_at DESC").Find(&publishModel).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(
		publishModel,
		mapper.PublisherModelToEntity,
		page,
		size,
		total,
	), nil
}

func (this *PublisherRepositoryImpl) GetAll() ([]*entity2.PublisherEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (this *PublisherRepositoryImpl) GetByID(id uint) (*entity2.PublisherEntity, error) {
	var publisher models.PublisherModel
	if result := this.db().Preload("Institution").Preload("Region").Preload("District").First(&publisher, id); result.Error != nil {
		return nil, infraError.Wrap(result.Error)
	}
	publisherEntity := mapper.PublisherModelToEntity(&publisher)
	return publisherEntity, nil
}

func (this *PublisherRepositoryImpl) GetByTin(tin string) (*entity2.PublisherEntity, error) {
	var publisher models.PublisherModel
	if result := this.db().Preload("Institution").Preload("Region").Preload("District").Where("tin = ?", tin).First(&publisher); result.Error != nil {
		return nil, infraError.Wrap(result.Error)
	}
	publisherEntity := mapper.PublisherModelToEntity(&publisher)
	return publisherEntity, nil
}

func (this *PublisherRepositoryImpl) Update(id uint, entity *entity2.PublisherEntity) error {
	entityModel := mapper.PublisherEntityToModel(entity)
	if result := this.db().Model(&models.PublisherModel{}).
		Where("id = ?", id).
		Select("tin", "name", "short_name", "logo", "phone_number", "fax_number", "email", "website", "address", "description", "is_active", "region_id", "district_id").
		Updates(entityModel); result.Error != nil {
		return infraError.Wrap(result.Error)
	}
	return nil
}

func (this *PublisherRepositoryImpl) Create(entity *entity2.PublisherEntity) error {
	publisherMap := mapper.PublisherEntityToModel(entity)
	if result := this.db().Create(&publisherMap); result.Error != nil {
		return infraError.Wrap(result.Error)
	}
	return nil
}

func (this *PublisherRepositoryImpl) Delete(id uint) error {
	var publisher models.PublisherModel
	if err := this.db().First(&publisher, id).Error; err != nil {
		return err
	}
	if err := this.db().Unscoped().Delete(&publisher).Error; err != nil {
		return err
	}
	return nil
}
