package usecase

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationDetailUsecase struct {
	repository    repository.ApplicationRepository
	referenceRepo repository.ArticleReferenceRepository
}

// @inject
func NewApplicationDetailUsecase(
	repository repository.ApplicationRepository,
	referenceRepo repository.ArticleReferenceRepository,
) *ApplicationDetailUsecase {
	return &ApplicationDetailUsecase{
		repository:    repository,
		referenceRepo: referenceRepo,
	}
}

func (this *ApplicationDetailUsecase) Execute(requester *entity.UserBasicEntity, applicationId uint) (*entity.ApplicationResponseEntity, error) {
	application, err := this.repository.GetByIDWithRelations(applicationId)
	if err != nil {
		return nil, err
	}

	references, err := this.referenceRepo.GetListByArticleID(application.ArticleID)
	if err != nil {
		return nil, err
	}
	application.Article.References = references

	response := mapper.ApplicationEntityToResponse(application)

	var requesterID uint
	if requester != nil {
		requesterID = requester.ID
	}
	permissionusecases.RedactApplicationContacts(response, requesterID, permissionusecases.IsAdmin(requester))

	return response, nil
}
