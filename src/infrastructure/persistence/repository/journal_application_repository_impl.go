package repository

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JournalApplicationRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalApplicationRepository(repository *BaseRepository) repository.JournalApplicationRepository {
	return &JournalApplicationRepositoryImpl{BaseRepository: repository}
}

func (this *JournalApplicationRepositoryImpl) Create(ent *entity.JournalCreateEntity) error {

	return this.db().Transaction(func(tx *gorm.DB) error {
		// journal create
		journalModel := mapper.JournalCreateEntityToModel(ent)
		if err := tx.Create(&journalModel).Error; err != nil {
			return err
		}
		// application create
		application := models.JournalApplicationModel{JournalID: &journalModel.ID}
		if err := tx.Create(&application).Error; err != nil {
			return err
		}

		// indexes create
		if len(ent.Indexes) > 0 {
			indexes := mapper.JournalIndexingToModel(journalModel.ID, ent.Indexes)
			if err := tx.Create(&indexes).Error; err != nil {
				return _errors.Wrap(err)
			}
		}

		// append study fields
		if err := tx.Model(journalModel).Association("StudyFields").Append(journalModel.StudyFields); err != nil {
			return err
		}

		// append languages
		if err := tx.Model(journalModel).Association("Languages").Append(journalModel.Languages); err != nil {
			return err
		}

		return nil

	})

}

func (this *JournalApplicationRepositoryImpl) GetListByPaging(publisherID uint, page, size int, status *enum.Status) (*entity.PagingEntity[entity.JournalApplicationEntity], error) {
	var applications []*models.JournalApplicationModel

	query := this.db().Model(&models.JournalApplicationModel{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if publisherID != 0 {
		query = query.Joins("JOIN journals ON journals.id = journal_applications.journal_id").
			Where("journals.publisher_id = ?", publisherID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := query.
		Preload("Journal", func(q *gorm.DB) *gorm.DB {
			return q.
				Omit("Description", "Rule", "ArticlePublishConditions", "Address")
		}).
		Preload("Journal.Publisher").
		Limit(size).
		Offset((page - 1) * size).
		Order("created_at DESC").
		Find(&applications).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.PagingMapper(applications, mapper.JournalApplicationModelToEntity, page, size, total), nil
}

// ISSNConditions ISSN turini to'liq SQL shartiga bog'laydi.
//
// Jadvalda ustun nomi emas, BUTUN shart saqlanadi. Aks holda quyida
// condition + " = ?" birlashtirish kerak bo'lardi — ya'ni aynan biz
// yo'q qilayotgan naqsh.
//
// Handler ham issn_type ni tekshiradi (CheckISSNHandler); bu yerdagi
// tekshiruv undan mustaqil va so'rov qurilayotgan joyning o'zida turadi.
var ISSNConditions = map[string]string{
	"issn_paper":  "journals.issn_paper = ?",
	"issn_online": "journals.issn_online = ?",
}

func (this *JournalApplicationRepositoryImpl) GetByISSN(issn string, issnField string) (*entity.JournalApplicationEntity, error) {
	var application models.JournalApplicationModel

	condition, ok := ISSNConditions[issnField]
	if !ok {
		return nil, response.NewFailResponse(400, "invalid ISSN type")
	}

	if err := (this.db().
		Preload("Journal").
		Joins("JOIN journals ON journals.id = journal_applications.journal_id").
		Where(condition, issn).
		Order("created_at desc").
		First(&application)).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalApplicationModelToEntity(&application), nil
}

func (this *JournalApplicationRepositoryImpl) SetStatus(id uint, status enum.Status, rejectionReason *string) error {

	application := models.JournalApplicationModel{
		Status:          status,
		RejectionReason: rejectionReason,
		ReviewedAt:      time.Now(),
	}
	return this.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.JournalApplicationModel{}).Where("id = ?", id).Updates(application).Error; err != nil {
			return _errors.Wrap(err)
		}
		var journalID uint
		if err := tx.Model(&models.JournalApplicationModel{}).Where("id = ?", id).Select("journal_id").Scan(&journalID).Error; err != nil {
			return _errors.Wrap(err)
		}

		if status == enum.StatusAccepted {
			journal := models.JournalModel{
				Model:      gorm.Model{ID: journalID},
				IsActive:   true,
				IsAccepted: true,
			}
			if err := tx.Model(&models.JournalModel{}).Where("id = ?", journalID).Updates(&journal).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (this *JournalApplicationRepositoryImpl) GetByID(id uint) (*entity.JournalApplicationEntity, error) {
	var journalApplication models.JournalApplicationModel

	if err := this.db().
		Preload("Journal").
		Preload("Journal.Publisher").
		Preload("Journal.StudyFields").
		Preload("Journal.Indexes").
		Preload("Journal.Languages").
		Preload("Journal.Region").
		Preload("Journal.District").
		Preload("User").
		Where("id = ?", id).
		First(&journalApplication).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalApplicationModelToEntity(&journalApplication), nil
}
