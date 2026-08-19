package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticlesListUseCase struct {
	repository repository.PublishedArticleRepository
}

// @inject
func NewArticlesListUseCase(repository repository.PublishedArticleRepository) *ArticlesListUseCase {
	return &ArticlesListUseCase{repository: repository}
}

func (this *ArticlesListUseCase) Execute(
	fromYear, toYear *int, journal uint, languages []uint,
	studyfields []uint,
	tags []string,
	articleTitle, authorName, journalName, issnPaper, issnOnline *string,
	sort string,
	page, pageSize int,
	accessType int,
	hasDoi *bool,
) (*entity2.PagingEntity[entity2.ArticleInputEntity], error) {
	var accessTypeEnum enum.AccessType
	switch accessType {
	case int(enum.PublicAccessType):
		accessTypeEnum = enum.PublicAccessType
	case int(enum.PrivateAccessType):
		accessTypeEnum = enum.PrivateAccessType
	default:
		accessTypeEnum = 0
	}

	paging, err := this.repository.GetAll(fromYear, toYear, journal, languages, studyfields, tags, articleTitle, authorName, journalName, issnPaper, issnOnline, sort, page, pageSize, accessTypeEnum, hasDoi)

	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(page, pageSize, paging.Total, mapper.ArticleEntityListToInput(paging.Items)), nil

}
