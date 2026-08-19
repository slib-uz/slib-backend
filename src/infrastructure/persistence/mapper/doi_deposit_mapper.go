package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func DoiDepositEntityToModel(it *entity.DoiDepositEntity) *models.DoiDepositModel {
	return models.NewDoiDepositModel(it.ArticleID, it.BatchID, it.DOI, it.Status, it.Message, it.SubmissionID, it.RequestBody, it.ResponseBody)
}

func DoiDepositModelToEntity(it *models.DoiDepositModel) *entity.DoiDepositEntity {
	return entity.NewDoiDepositEntity(it.ID, it.ArticleID, it.BatchID, it.DOI, it.Status, it.Message, it.SubmissionID, it.RequestBody, it.ResponseBody, it.CreatedAt, it.UpdatedAt)
}
