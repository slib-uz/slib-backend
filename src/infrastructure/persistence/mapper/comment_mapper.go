package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func CommentModelToEntity(model *models.CommentModel) *entity2.CommentEntity {

	var article *entity2.ArticleInputEntity
	var user *entity2.UserSharedEntity
	if model.Article != nil {
		article = ArticleModelToInput(model.Article)
	}

	if model.User != nil {
		user = UserModelToSharedEntity(model.User)
	}

	return entity2.NewCommentEntity(
		model.ID,
		model.ArticleID,
		article,
		model.UserID,
		user,
		model.Content,
		model.Rating,
		model.CreatedAt,
	)
}

func CommentEntityToModel(comment *entity2.CommentEntity) *models.CommentModel {

	return models.NewCommentModel(
		comment.ArticleID,
		comment.UserID,
		comment.Content,
		comment.Rating,
	)
}
