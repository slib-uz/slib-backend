package repository

import (
	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type NotificationRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNotificationRepository(repository *BaseRepository) repository.NotificationRepository {
	return &NotificationRepositoryImpl{BaseRepository: repository}
}

func (this *NotificationRepositoryImpl) Create(notification *entity2.NotificationEntity) (uint, error) {
	var _model = mapper.NotificationEntityToModel(notification)
	if err := this.db().Create(_model).Error; err != nil {
		return 0, infraError.Wrap(err)
	}
	return _model.ID, nil
}

func (this *NotificationRepositoryImpl) BulkCreate(notifications []*entity2.NotificationEntity) ([]uint, error) {
	var _models = mapper.NotificationEntityListToModelList(notifications)

	if err := this.db().Create(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	var ids = make([]uint, len(_models))
	for i, model := range _models {
		ids[i] = model.ID
	}
	return ids, nil
}

func (this *NotificationRepositoryImpl) GetByID(id, userID uint) (*entity2.NotificationEntity, error) {
	var _model models.NotificationModel

	if err := this.getWithReadStatus(this.db(), userID).First(&_model, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.NotificationModelToEntity(&_model), nil
}

func (this *NotificationRepositoryImpl) GetByUserIDAndBroadcast(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.NotificationEntity], error) {
	var _models []*models.NotificationModel

	var total int64

	if err := this.db().Model(&models.NotificationModel{}).Where("user_id = ? OR is_broadcast = true", userID).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := this.getWithReadStatus(this.db(), userID).
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Where("n.user_id = ? OR n.is_broadcast = ?", userID, true).
		Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(
		_models,
		mapper.NotificationModelToEntity,
		page,
		pageSize,
		total,
	), nil
}

func (this *NotificationRepositoryImpl) Read(id, userID uint) error {
	var notification models.NotificationModel

	if err := this.db().Select("id, user_id, is_broadcast, is_user_read").First(&notification, id).Error; err != nil {
		return infraError.Wrap(err)
	}

	if notification.IsBroadcast {
		var readModel = models.NewNotificationReadModel(notification.ID, userID)
		return this.db().Create(&readModel).Error
	}

	if notification.IsUserRead {
		return nil // Already read
	}
	return this.db().Model(&notification).Update("is_user_read", true).Error

}

func (this *NotificationRepositoryImpl) GetUnreadNotificationsByUserID(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.NotificationEntity], error) {
	var _models []*models.NotificationModel

	query := this.getWithReadStatus(this.db(), userID).
		Where(`
			(n.is_broadcast = false AND n.is_user_read = false)
			OR
			(n.is_broadcast = true AND nr.id IS NULL)
		`)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(
		_models,
		mapper.NotificationModelToEntity,
		page,
		pageSize,
		total,
	), nil
}

func (this *NotificationRepositoryImpl) GetUnreadCount(userID uint) (int64, error) {
	var count int64

	err := this.db().Table("notification_models n").
		Joins(`
            LEFT JOIN notification_read_models nr
                ON nr.notification_id = n.id AND nr.user_id = ?
        `, userID).
		Where(`
            (n.user_id = ? OR n.is_broadcast = true)
            AND (
                (n.is_broadcast = false AND n.is_user_read = false)
                OR
                (n.is_broadcast = true AND nr.id IS NULL)
            )
        `, userID).
		Count(&count).Error
	return count, err
}

func (this *NotificationRepositoryImpl) getWithReadStatus(query *gorm.DB, userID uint) *gorm.DB {
	return query.
		Table("notification_models n").
		Select(`
			n.*,
			CASE
				WHEN n.is_broadcast = false THEN n.is_user_read
				ELSE CASE WHEN nr.id IS NULL THEN false ELSE true END
			END AS is_read
		`).
		Joins(`LEFT JOIN notification_read_models nr ON nr.notification_id = n.id AND nr.user_id = ?`, userID).
		Order("n.created_at DESC")
}
