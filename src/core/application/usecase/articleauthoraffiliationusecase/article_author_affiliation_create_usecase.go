package articleauthoraffiliationusecase

import (
	"slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/ports/repository"
)

type ArticleAuthorAffiliationCreateUseCase struct {
	repository             repository.ArticleAuthorAffiliationRepository
	organizationRepository repository.OrganizationRepository
}

// @inject
func NewArticleAuthorAffiliationCreateUseCase(repository repository.ArticleAuthorAffiliationRepository, organizationRepository repository.OrganizationRepository) *ArticleAuthorAffiliationCreateUseCase {
	return &ArticleAuthorAffiliationCreateUseCase{repository: repository, organizationRepository: organizationRepository}
}

func (this *ArticleAuthorAffiliationCreateUseCase) Execute(articleAuthorAffiliation *entity.ArticleAuthorAffiliationEntity) (*entity.ArticleAuthorAffiliationEntity, error) {
	if articleAuthorAffiliation.OrganizationID != nil {
		organization, err := this.organizationRepository.GetByID(*articleAuthorAffiliation.OrganizationID)
		if err != nil {
			return nil, err
		}

		articleAuthorAffiliation.OrganizationName = organization.Name
		if organization.Tin != nil {
			articleAuthorAffiliation.OrganizationTin = *organization.Tin
		}
	}

	return this.repository.Create(articleAuthorAffiliation)

	//articleAuthorAffiliationEntity, err := this.repository.Create(articleAuthorAffiliation)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating article author affiliation: %w", err)
	//}
	//return articleAuthorAffiliationEntity, nil
}
