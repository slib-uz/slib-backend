package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func OrganizationModelToEntity(data *models.OrganizationModel) *entity.OrganizationEntity {
	tin := &data.Tin
	e := entity.NewOrganizationEntity(data.ID, data.Name, tin, data.Address, data.SoatoID)
	if data.Soato != nil {
		e.Soato = SoatoModelToEntity(data.Soato)
	}
	return e
}

func OrganizationEntityToModel(data *entity.OrganizationEntity) *models.OrganizationModel {
	tin := ""
	if data.Tin != nil {
		tin = *data.Tin
	}
	return models.NewOrganizationModel(tin, data.Name, data.Address, data.SoatoID)
}
