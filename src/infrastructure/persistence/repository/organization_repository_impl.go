package repository

import (
	entity2 "slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type OrganizationRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewOrganizationRepository(repository *BaseRepository) repository.OrganizationRepository {
	return &OrganizationRepositoryImpl{BaseRepository: repository}
}

func (this *OrganizationRepositoryImpl) GetByTin(tin string) (*entity2.OrganizationEntity, error) {
	var organization models.OrganizationModel
	if err := this.db().Where("tin = ?", tin).First(&organization).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.OrganizationModelToEntity(&organization), nil
}

func (this *OrganizationRepositoryImpl) GetByID(id uint) (*entity2.OrganizationEntity, error) {
	var organization models.OrganizationModel
	if err := this.db().Preload("Soato").First(&organization, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.OrganizationModelToEntity(&organization), nil
}

func (this *OrganizationRepositoryImpl) GetList(page, size int, tin, name, address *string) (*entity2.PagingEntity[entity2.OrganizationEntity], error) {
	var organizations []*models.OrganizationModel

	query := this.db().Model(&models.OrganizationModel{})

	if tin != nil && *tin != "" {
		query = query.Where("tin ILIKE ?", "%"+*tin+"%")
	}
	if name != nil && *name != "" {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}
	if address != nil && *address != "" {
		query = query.Where("address ILIKE ?", "%"+*address+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Limit(size).Offset((page - 1) * size).Order("created_at DESC").Find(&organizations).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	entities := make([]*entity2.OrganizationEntity, len(organizations))
	for i, org := range organizations {
		entities[i] = mapper.OrganizationModelToEntity(org)
	}

	return entity2.NewPagingEntity(page, size, total, entities), nil
}

func (this *OrganizationRepositoryImpl) Create(e *entity2.OrganizationEntity) error {
	model := mapper.OrganizationEntityToModel(e)
	if err := this.db().Create(model).Error; err != nil {
		return infraError.Wrap(err)
	}
	*e = *mapper.OrganizationModelToEntity(model)
	return nil
}

func (this *OrganizationRepositoryImpl) Update(id uint, e *entity2.OrganizationEntity) error {
	var model models.OrganizationModel
	if err := this.db().First(&model, id).Error; err != nil {
		return infraError.Wrap(err)
	}

	newModel := mapper.OrganizationEntityToModel(e)
	model.Tin = newModel.Tin
	model.Name = newModel.Name
	model.Address = newModel.Address
	model.SoatoID = newModel.SoatoID

	if err := this.db().Save(&model).Error; err != nil {
		return infraError.Wrap(err)
	}
	*e = *mapper.OrganizationModelToEntity(&model)
	return nil
}

func (this *OrganizationRepositoryImpl) Delete(id uint) error {
	if err := this.db().Delete(&models.OrganizationModel{}, id).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}
