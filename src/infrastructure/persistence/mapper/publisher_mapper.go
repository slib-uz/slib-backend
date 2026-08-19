package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func PublisherEntityToModel(entity *entity.PublisherEntity) *models.PublisherModel {
	model := models.NewPublisherModel(
		entity.Tin,
		entity.Name,
		entity.ShortName,
		entity.Logo,
		entity.PhoneNumber,
		entity.FaxNumber,
		entity.Email,
		entity.Website,
		entity.Address,
		entity.Description,
		entity.IsActive,
	)
	model.InstitutionID = entity.InstitutionID
	model.RegionID = entity.RegionID
	model.DistrictID = entity.DistrictID
	return model
}

func PublisherModelToEntity(model *models.PublisherModel) *entity.PublisherEntity {
	publisher := entity.NewPublisherEntity(
		model.ID,
		model.Tin,
		model.Name,
		model.ShortName,
		model.Logo,
		model.PhoneNumber,
		model.FaxNumber,
		model.Email,
		model.Website,
		model.Address,
		model.Description,
		model.IsActive,
	)
	if model.Institution != nil {
		publisher.Institution = InstitutionModelToEntity(model.Institution)
	}
	publisher.InstitutionID = model.InstitutionID
	publisher.RegionID, publisher.DistrictID, publisher.Region, publisher.District = MapRegionDistrictFromModel(
		model.RegionID,
		model.DistrictID,
		model.Region,
		model.District,
	)
	return publisher
}

func PublisherUpdateModel(existing, update *models.PublisherModel) *models.PublisherModel {
	if update.Tin != "" && update.Tin != existing.Tin {
		existing.Tin = update.Tin
	}
	if update.Name != nil && update.Name != existing.Name {
		existing.Name = update.Name
	}
	if update.ShortName != nil {
		existing.ShortName = update.ShortName
	}
	if update.Logo != nil {
		existing.Logo = update.Logo
	}
	if update.PhoneNumber != nil {
		existing.PhoneNumber = update.PhoneNumber
	}
	if update.FaxNumber != nil {
		existing.FaxNumber = update.FaxNumber
	}
	if update.Email != nil {
		existing.Email = update.Email
	}
	if update.Website != nil {
		existing.Website = update.Website
	}
	if update.Address != nil {
		existing.Address = update.Address
	}
	if update.Description != nil {
		existing.Description = update.Description
	}
	//if existing.Status == 0 {
	//	existing.Status = update.Status
	//}
	return existing

}
