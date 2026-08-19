package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func ReportListModelToEntity(data []*models.ReportModel) []*entity.ReportEntity {
	var result = make([]*entity.ReportEntity, len(data))

	for i, item := range data {
		result[i] = ReportModelToEntity(item)
	}

	return result
}

func ReportEntityToModel(data *entity.ReportEntity) *models.ReportModel {
	return models.NewReportModel(data.Reason, data.TargetID, data.TargetType, data.ReporterID)
}

func ReportModelToEntity(data *models.ReportModel) *entity.ReportEntity {
	e := entity.NewReportEntity(data.ID, data.Reason, data.TargetID, data.Target, data.TargetType, data.ReporterID, ReporterModelToEntity(data.Reporter), data.CreatedAt)
	if data.Target != nil {
		switch t := data.Target.(type) {
		case *models.ArticleModel:
			e.Target = ArticleModelToEntity(t)
		case *models.JournalModel:
			e.Target = JournalModelToEntity(t)
		default:
			return nil
		}
	}
	return e
}
