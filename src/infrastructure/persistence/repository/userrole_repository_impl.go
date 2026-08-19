package repository

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type UserRoleRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewUserRoleRepository(repository *BaseRepository) repository.UserRoleRepository {
	return &UserRoleRepositoryImpl{BaseRepository: repository}
}

func (this *UserRoleRepositoryImpl) GetByUserID(userID uint) ([]*entity.UserRoleEntity, error) {
	var userRoles []*models.RoleModel
	if result := this.db().
		Preload("Journal").
		Preload("Publisher").
		Preload("Institution").
		Where("user_id = ?", userID).
		Find(&userRoles); result.Error != nil {
		return nil, result.Error
	}

	userRoleEntities := make([]*entity.UserRoleEntity, len(userRoles))
	for i, role := range userRoles {
		userRoleEntities[i] = mapper.UserRoleModelToEntity(role)
	}

	return userRoleEntities, nil
}

func (this *UserRoleRepositoryImpl) Create(userRole *entity.UserRoleEntity) error {
	roleModel := mapper.UserRoleEntityToModel(userRole)
	if result := this.db().Create(&roleModel); result.Error != nil {
		return result.Error
	}
	return nil
}

func (this *UserRoleRepositoryImpl) Exists(userID uint, publisherID, journalID, institutionID *uint) (bool, error) {
	var count int64
	if err := this.db().Model(&models.RoleModel{}).
		Where("user_id = ? AND publisher_id = ? AND journal_id = ? AND institution_id = ?", userID, publisherID, journalID, institutionID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func (this *UserRoleRepositoryImpl) GetByPublisherID(publisherID uint) ([]*entity.UserRoleEntity, error) {
	var userRoles []*models.RoleModel

	if err := this.db().
		Preload("User").
		Preload("Publisher").
		Where("publisher_id = ? and role = ?", publisherID, enum.RolePublisherAdmin).
		Find(&userRoles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	userRoleEntities := make([]*entity.UserRoleEntity, len(userRoles))
	for i, role := range userRoles {
		userRoleEntities[i] = mapper.UserRoleModelToEntity(role)
	}
	return userRoleEntities, nil
}

func (this *UserRoleRepositoryImpl) GetByInstitutionID(institutionID uint) ([]*entity.UserRoleEntity, error) {
	var userRoles []*models.RoleModel

	if err := this.db().
		Preload("User").
		Preload("Institution").
		Where("institution_id = ? and role = ?", institutionID, enum.RoleInstitutionAdmin).
		Find(&userRoles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	userRoleEntities := make([]*entity.UserRoleEntity, len(userRoles))
	for i, role := range userRoles {
		userRoleEntities[i] = mapper.UserRoleModelToEntity(role)
	}
	return userRoleEntities, nil
}

func (this *UserRoleRepositoryImpl) Delete(id uint) error {
	err := this.db().Delete(&models.RoleModel{}, id).Error
	return infraError.Wrap(err)

}

func (this *UserRoleRepositoryImpl) GetByJournalID(journalID uint) ([]*entity.UserRoleEntity, error) {
	var _roles []*models.RoleModel

	if err := this.db().
		Where("journal_id = ?", journalID).
		Preload("User").
		Find(&_roles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	var roles = make([]*entity.UserRoleEntity, len(_roles))
	for i, _role := range _roles {
		roles[i] = mapper.UserRoleModelToEntity(_role)
	}
	return roles, nil
}

func (this *UserRoleRepositoryImpl) FindByJournalID(journalID uint) ([]*entity.UserRoleEntity, error) {
	var _roles []*models.RoleModel

	if err := this.db().
		Where("journal_id = ?", journalID).
		Find(&_roles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	var roles = make([]*entity.UserRoleEntity, len(_roles))
	for i, _role := range _roles {
		roles[i] = mapper.UserRoleModelToEntity(_role)
	}
	return roles, nil
}

func (this *UserRoleRepositoryImpl) GetJournalMemberIds(journalID uint) ([]uint, error) {
	var ids []uint
	if err := this.db().
		Model(&models.RoleModel{}).
		Where("journal_id = ?", journalID).
		Pluck("user_id", &ids).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return ids, nil
}

type journalMemberLastOnlineScan struct {
	RoleID       uint          `gorm:"column:role_id"`
	UserID       uint          `gorm:"column:user_id"`
	FullName     string        `gorm:"column:full_name"`
	ScienceID    string        `gorm:"column:science_id"`
	PhoneNumber  string        `gorm:"column:phone_number"`
	Role         enum.UserRole `gorm:"column:role"`
	LastOnlineAt *time.Time    `gorm:"column:last_online_at"`
}

func (this *UserRoleRepositoryImpl) GetByJournalIDOrderByLastOnline(journalID uint) ([]*entity.JournalMemberLastOnlineEntity, error) {
	var rows []journalMemberLastOnlineScan

	query := `
		SELECT
			r.id AS role_id,
			r.user_id,
			r.role,
			u.full_name,
			u.science_id,
			u.phone_number,
			up.last_online_at
		FROM roles r
		JOIN users u ON u.id = r.user_id AND u.deleted_at IS NULL
		LEFT JOIN user_profiles up ON up.user_id = r.user_id AND up.deleted_at IS NULL
		WHERE r.journal_id = ? AND r.deleted_at IS NULL
		ORDER BY up.last_online_at DESC NULLS LAST
	`

	if err := this.db().Raw(query, journalID).Scan(&rows).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	members := make([]*entity.JournalMemberLastOnlineEntity, len(rows))
	for i, row := range rows {
		members[i] = entity.NewJournalMemberLastOnlineEntity(
			row.RoleID,
			row.UserID,
			row.FullName,
			row.ScienceID,
			row.PhoneNumber,
			row.Role,
			row.LastOnlineAt,
		)
	}

	return members, nil
}
