package repository

import (
	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/db"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type StudyFieldRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewStudyFieldRepository(database *db.Database) repository.StudyFieldRepository {
	return &StudyFieldRepositoryImpl{database: database}
}

func (this *StudyFieldRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *StudyFieldRepositoryImpl) GetAll(search string) ([]*entity2.StudyFieldEntity, error) {

	var studyFieldsModels []models.StudyFieldModel

	query := this.db().Model(&models.StudyFieldModel{}).Where("parent_id IS NULL")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(`name->>'uz' ILIKE ? OR name->>'ru' ILIKE ? OR name->>'en' ILIKE ?`, searchPattern, searchPattern, searchPattern)
	}

	if result := query.Preload("Children").Find(&studyFieldsModels); result.Error != nil {
		return nil, result.Error
	}

	var studyFields = make([]*entity2.StudyFieldEntity, len(studyFieldsModels))
	for i, studyFieldModel := range studyFieldsModels {
		studyFields[i] = mapper.StudyFieldModelToEntity(&studyFieldModel)
	}

	return studyFields, nil

}

func (this *StudyFieldRepositoryImpl) ExistingIds(ids []uint) ([]uint, error) {
	var existingIds []uint
	if err := this.db().Model(&models.StudyFieldModel{}).Where("id IN ?", ids).Pluck("id", &existingIds).Error; err != nil {
		return nil, err
	}
	return existingIds, nil
}

func (this *StudyFieldRepositoryImpl) Create(studyField *entity2.StudyFieldEntity) error {
	var _model = mapper.StudyFieldEntityToModel(studyField)
	return this.db().Create(&_model).Error
}

func (this *StudyFieldRepositoryImpl) Delete(id uint) error {
	return this.db().Delete(&models.StudyFieldModel{}, id).Error
}

func (this *StudyFieldRepositoryImpl) Update(studyField *entity2.StudyFieldEntity) error {
	var _model = mapper.StudyFieldEntityToModel(studyField)

	return this.db().Updates(_model).Error
}

func (this *StudyFieldRepositoryImpl) GetByPaging(page, pageSize int, search string) (*entity2.PagingEntity[entity2.StudyFieldEntity], error) {
	var _models []models.StudyFieldModel

	query := this.db().Model(&models.StudyFieldModel{}).Where("parent_id IS NULL")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(`name->>'uz' ILIKE ? OR name->>'ru' ILIKE ? OR name->>'en' ILIKE ?`, searchPattern, searchPattern, searchPattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Preload("Children").Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, err
	}

	var items = make([]*entity2.StudyFieldEntity, len(_models))
	for i, model := range _models {
		items[i] = mapper.StudyFieldModelToEntity(&model)
	}
	return entity2.NewPagingEntity(page, pageSize, total, items), nil
}

func (this *StudyFieldRepositoryImpl) GetByJournalID(journalID uint) ([]*entity2.StudyFieldEntity, error) {
	// Получаем все study_field_id, которые связаны с journal
	var journalStudyFieldIDs []uint
	if err := this.db().Table("journal_many2many_study_fields").
		Where("journal_model_id = ?", journalID).
		Pluck("study_field_model_id", &journalStudyFieldIDs).Error; err != nil {
		return nil, err
	}

	if len(journalStudyFieldIDs) == 0 {
		return []*entity2.StudyFieldEntity{}, nil
	}

	// Получаем parent_id из связанных studyfield'ов напрямую из БД
	var parentIDs []uint
	if err := this.db().Model(&models.StudyFieldModel{}).
		Where("id IN ? AND parent_id IS NOT NULL", journalStudyFieldIDs).
		Pluck("parent_id", &parentIDs).Error; err != nil {
		return nil, err
	}

	// Создаем map для уникальности parent_id
	parentIDsMap := make(map[uint]bool)
	for _, id := range parentIDs {
		parentIDsMap[id] = true
	}

	// Создаем map для быстрого доступа к связанным с journal studyfield'ам
	journalStudyFieldMap := make(map[uint]bool)
	for _, id := range journalStudyFieldIDs {
		journalStudyFieldMap[id] = true
	}

	// Получаем все нужные studyfield'ы одним запросом:
	// 1. Parents (даже если они не связаны с journal)
	// 2. Children, которые связаны с journal
	var parentIDsList []uint
	for id := range parentIDsMap {
		parentIDsList = append(parentIDsList, id)
	}

	var allStudyFields []models.StudyFieldModel
	if len(parentIDsList) > 0 {
		// Объединяем parentIDsList и journalStudyFieldIDs в один список для запроса
		allIDs := make(map[uint]bool)
		for _, id := range parentIDsList {
			allIDs[id] = true
		}
		for _, id := range journalStudyFieldIDs {
			allIDs[id] = true
		}
		allIDsList := make([]uint, 0, len(allIDs))
		for id := range allIDs {
			allIDsList = append(allIDsList, id)
		}

		// Получаем parents и children, которые связаны с journal
		if err := this.db().
			Where("id IN ?", allIDsList).
			Preload("Children", "id IN ?", journalStudyFieldIDs).
			Find(&allStudyFields).Error; err != nil {
			return nil, err
		}
	} else {
		// Если нет parents, возвращаем только связанные studyfield'ы (root level)
		if err := this.db().
			Where("id IN ? AND parent_id IS NULL", journalStudyFieldIDs).
			Find(&allStudyFields).Error; err != nil {
			return nil, err
		}
	}

	// Создаем map для всех studyfield'ов
	studyFieldMap := make(map[uint]*entity2.StudyFieldEntity)

	// Конвертируем в entities
	for i := range allStudyFields {
		model := &allStudyFields[i]
		studyFieldEntity := mapper.StudyFieldModelToEntity(model)

		// Если это parent, фильтруем children - оставляем только те, что связаны с journal
		if len(model.Children) > 0 {
			var filteredChildren []*entity2.StudyFieldEntity
			for _, child := range model.Children {
				if journalStudyFieldMap[child.ID] {
					childEntity := mapper.StudyFieldModelToEntity(&child)
					filteredChildren = append(filteredChildren, childEntity)
				}
			}
			studyFieldEntity.Children = filteredChildren
		}

		studyFieldMap[studyFieldEntity.ID] = studyFieldEntity
	}

	// Формируем результат: только root level (parent_id IS NULL)
	var result []*entity2.StudyFieldEntity
	addedMap := make(map[uint]bool)

	for _, studyFieldEntity := range studyFieldMap {
		if studyFieldEntity.ParentID == nil {
			if !addedMap[studyFieldEntity.ID] {
				result = append(result, studyFieldEntity)
				addedMap[studyFieldEntity.ID] = true
			}
		}
	}

	return result, nil
}
