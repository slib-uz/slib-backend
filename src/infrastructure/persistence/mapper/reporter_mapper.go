package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ReporterModelToEntity(data *models.UserModel) *entity.ReporterEntity {
	return entity.NewReporterEntity(data.ID, data.FullName, data.ScienceID)
}
