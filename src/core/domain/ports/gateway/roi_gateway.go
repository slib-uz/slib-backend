package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type ROIGateway interface {
	Publish(data *entity.ArticleEntity, sourceLink string) (*entity.ROIDataEntity, error)
	UpdateArticle(article *entity.ROIArticleUpdateEntity) error
}
