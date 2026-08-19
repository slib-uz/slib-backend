package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/session"
)

type UserRepository interface {
	GetById(id uint) (*entity.UserEntity, error)
	FindByPhone(phone string) (*entity.UserEntity, error, bool)
	Save(user *entity.UserEntity) (*entity.UserEntity, error)
	Update(id uint, user *entity.UserEntity) (*entity.UserEntity, error)
	GetByScienceId(scienceId string) (*entity.UserEntity, error)
	GetGenderStatistics() (*entity.UserGenderStatisticsEntity, error)
	GetAgeStatistics() (*entity.UserAgeStatisticsEntity, error)
	GetArticleStatisticsByYear(userID uint, year int) (*entity.UserArticleStatisticsByYearEntity, error)
	GetArticleStatisticsByYearRange(userID uint, fromYear, toYear int) (*entity.UserArticleStatisticsByYearRangeEntity, error)
	GetUserActivityStartYear(userID uint) (int, error)
	GetAll(page, pageSize int, search string) (*entity.PagingEntity[entity.UserEntity], error)
	GetDetailByID(id uint) (*entity.UserDetailEntity, error)
	GetByOrcid(orcid string) (*entity.UserEntity, error)
	UpdateOrcid(userID uint, orcid string) error
	GetByPhoneNumber(phoneNumber string) (*entity.UserEntity, error)
	CreateByPhoneNumber(tx session.Tx, user *entity.UserEntity) (uint, error)
}
