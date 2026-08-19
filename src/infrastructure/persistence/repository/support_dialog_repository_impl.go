package repository

import (
	"fmt"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

// SupportDialogSortFields — qo'llab-quvvatlash dialoglari uchun ruxsat etilgan
// tartiblash maydonlari. Uchala metod ham SupportDialogModel ni so'raydi,
// TableName() "support_dialogs" qaytaradi.
var SupportDialogSortFields = sorting.New("support_dialogs.created_at DESC", map[string]string{
	"created_at": "support_dialogs.created_at",
})

type SupportDialogRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewSupportDialogRepository(baseRepository *BaseRepository) repository.SupportDialogRepository {
	return &SupportDialogRepositoryImpl{BaseRepository: baseRepository}
}

func (this SupportDialogRepositoryImpl) CreateAnswer(supportDialog *entity2.SupportDialogEntity) error {
	var _model = mapper.SupportDialogEntityToModel(supportDialog)
	var _chat *models.UserModel

	if err := this.db().Where("id = ?", supportDialog.ChatID).First(&_chat).Error; err != nil {
		return infraError.Wrap(err)
	}

	if err := this.db().Model(&models.SupportDialogModel{}).Create(_model).Error; err != nil {
		return err
	}
	return nil
}

func (this SupportDialogRepositoryImpl) CreateQuestion(supportDialog *entity2.SupportDialogEntity) error {
	var _model = mapper.SupportDialogEntityToModel(supportDialog)
	if err := this.db().Model(&models.SupportDialogModel{}).Create(_model).Error; err != nil {
		return err
	}
	return nil
}

func (this SupportDialogRepositoryImpl) GetByPaging(page, pageSize int, ordering string, userID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error) {
	orderExpr, err := SupportDialogSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.SupportDialogModel

	countQuery := this.db().Model(&models.SupportDialogModel{}).Where("chat_id = ?", userID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query := this.db().Model(&models.SupportDialogModel{}).Where("chat_id = ?", userID).Preload("Owner")

	query = query.Order(orderExpr)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, fmt.Errorf("error fetching support dialog")
	}

	supportDialog := mapper.SupportDialogListModelToEntity(_models)

	return entity2.NewPagingEntity(page, pageSize, total, supportDialog), nil
}

func (this SupportDialogRepositoryImpl) GetByChatID(page, pageSize int, ordering string, chatID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error) {
	orderExpr, err := SupportDialogSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.SupportDialogModel

	countQuery := this.db().Model(&models.SupportDialogModel{}).Where("chat_id = ?", chatID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query := this.db().Model(&models.SupportDialogModel{}).Where("chat_id = ?", chatID).Preload("Owner")

	query = query.Order(orderExpr)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, err
	}

	supportDialog := mapper.SupportDialogListModelToEntity(_models)

	return entity2.NewPagingEntity(page, pageSize, total, supportDialog), nil
}

func (this SupportDialogRepositoryImpl) GetChatsByPaging(page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.ChatEntity], error) {
	orderExpr, err := SupportDialogSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.SupportDialogModel

	if err := this.db().
		Model(&models.SupportDialogModel{}).
		Select("COUNT(DISTINCT chat_id)").
		Count(&total).Error; err != nil {
		return nil, err
	}

	subQuery := this.db().
		Table("support_dialogs").
		Select("chat_id, MAX(created_at) AS last_created_at").
		Where("message_type = ? AND deleted_at IS NULL", enum.Question).
		Group("chat_id")

	query := this.db().
		Model(&models.SupportDialogModel{}).
		Joins("JOIN (?) AS last_msg ON support_dialogs.chat_id = last_msg.chat_id AND support_dialogs.created_at = last_msg.last_created_at", subQuery).
		Preload("Owner")

	query = query.Order(orderExpr)

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&_models).Error; err != nil {
		return nil, err
	}

	var chats []*entity2.ChatEntity
	for _, model := range _models {
		unreadCount, err := this.GetUnreadCount(model.ChatID)
		if err != nil {
			return nil, err
		}
		chats = append(chats, entity2.NewChatEntity(mapper.SupportDialogModelToEntity(model), unreadCount))
	}

	return entity2.NewPagingEntity(page, pageSize, total, chats), nil
}

func (this SupportDialogRepositoryImpl) GetUnreadCount(chatID uint) (int64, error) {
	var count int64
	if err := this.db().
		Model(&models.SupportDialogModel{}).
		Where("chat_id = ? AND is_read = ?", chatID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (this SupportDialogRepositoryImpl) MarkAsRead(chatID uint) error {
	return this.db().
		Model(&models.SupportDialogModel{}).
		Where("chat_id = ? AND is_read = ?", chatID, false).
		Update("is_read", true).Error
}
