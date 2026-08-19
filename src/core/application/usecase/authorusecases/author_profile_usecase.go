package authorusecases

import (
	"slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type AuthorProfileUseCase struct {
	authorRepo repository.AuthorRepository
}

// @inject
func NewAuthorProfileUseCase(authorRepo repository.AuthorRepository) *AuthorProfileUseCase {
	return &AuthorProfileUseCase{authorRepo: authorRepo}
}

func (this *AuthorProfileUseCase) Execute(scienceID string) (*entity.UserProfileEntity, error) {
	author, err := this.authorRepo.GetByScienceID(scienceID)
	if err != nil {
		return nil, err
	}

	var academicDegree enum2.AcademicDegreeCode
	if author.AcademicDegree != nil {
		academicDegree = enum2.AcademicDegreeCode(*author.AcademicDegree)
	}

	var city enum2.CityCode
	if author.User != nil {
		city = enum2.CityCode(author.User.City)
	}

	return entity.NewUserProfileEntity(
		0,
		0,
		author.FullName,
		author.ScienceID,
		"",
		nil,
		author.Photo,
		"",
		nil,
		author.Job,
		author.ResearchMetrics,
		&city,
		&academicDegree,
		author.AcademicTitle,
		author.ORCIDID,
		0,
		0,
	), nil
}
