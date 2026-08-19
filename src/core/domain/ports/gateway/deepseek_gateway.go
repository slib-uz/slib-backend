package gateway

import "slib.uz/src/core/domain/entity"

type DeepSeekGateway interface {
	ExtractArticleMetadata(articleText string, studyFields []entity.StudyFieldCatalogItem, langs []string) (*entity.ArticleMetadataExtraction, error)
}
