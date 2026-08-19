package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type ScienceIDGateway interface {
	GetAuthorByScienceID(scienceID string) (*entity.AuthorEntity, error)
	GetUserScientificDataByScienceID(scienceID string) (*entity.AcademicDegreeEntity, *entity.AcademicTitleEntity, error)
	GetUserByScienceID(scienceID string) (*entity.UserEntity, error)
}
