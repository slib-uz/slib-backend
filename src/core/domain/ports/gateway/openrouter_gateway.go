package gateway

import "slib.uz/src/core/domain/entity"

type OpenRouterGateway interface {
	ExtractArticleMetadata(articleText string, studyFields []entity.StudyFieldCatalogItem, langs []string) (*entity.ArticleMetadataExtraction, error)
}
