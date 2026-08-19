package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"

	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ArticleSubmissionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewArticleSubmissionRepository(baseRepository *BaseRepository) repository.ArticleSubmissionRepository {
	return &ArticleSubmissionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *ArticleSubmissionRepositoryImpl) Create(
	article *entity2.ArticleCreateEntity,
	authorAffiliations []*entity2.ArticleAuthorAffiliationEntity,
	updateAuthorAffiliationIDs []uint,
	userId uint,
	technicalReviewDeadline time.Time) (*entity2.ApplicationEntity, error) {

	var journal models.JournalModel
	if err := this.db().First(&journal, article.JournalID).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	articleModel := mapper.ArticleCreateEntityToModel(article, journal.AccessType)
	var articleApplicationModel models.ArticleApplicationModel

	tr := this.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Session(&gorm.Session{FullSaveAssociations: true}).
			Omit("CoAuthors", "StudyFields", "Tags").
			Create(articleModel).Error; err != nil {
			return err
		}

		if err := tx.Model(articleModel).
			Association("CoAuthors").
			Append(articleModel.CoAuthors); err != nil {
			return err
		}
		if err := tx.Model(articleModel).
			Association("StudyFields").
			Append(articleModel.StudyFields); err != nil {
			return err
		}

		if err := tx.Model(articleModel).
			Association("Tags").
			Append(articleModel.Tags); err != nil {
			return err
		}

		// Create References
		if len(article.References) > 0 {
			var references []*models.ReferenceModel
			for _, refName := range article.References {
				references = append(references, models.NewReferenceModel(refName, articleModel.ID))
			}
			if err := tx.Create(&references).Error; err != nil {
				return err
			}
		}

		// Serialize number allocation per journal (unique index includes soft-deleted rows).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&models.JournalModel{}, articleModel.JournalID).Error; err != nil {
			return err
		}

		number, err := this.generateNumber(tx, articleModel.JournalID)
		if err != nil {
			return err
		}

		articleApplicationModel = *models.NewArticleApplicationModel(number, articleModel.ID, articleModel.JournalID, userId)

		if err := tx.Create(&articleApplicationModel).Error; err != nil {
			return err
		}

		pendingStage := models.NewReviewStageModel(0, articleApplicationModel.ID, enum2.TechnicalReviewStage, enum2.StatusPending, nil, technicalReviewDeadline, nil)
		if err := tx.Create(&pendingStage).Error; err != nil {
			return err
		}

		_models := mapper.ArticleAuthorAffiliationListEntityToModel(authorAffiliations)

		for _, model := range _models {
			model.ArticleID = &articleModel.ID
		}

		if len(_models) > 0 {
			if err := tx.Create(_models).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(models.ArticleAuthorAffiliationModel{}).
			Where("id IN ?", updateAuthorAffiliationIDs).
			UpdateColumn("article_id", articleModel.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if tr != nil {
		return nil, tr
	}

	return mapper.ApplicationModelToEntity(&articleApplicationModel), nil
}

func (this *ArticleSubmissionRepositoryImpl) generateNumber(tx *gorm.DB, journalID uint) (string, error) {
	year := time.Now().Year()
	var last models.ArticleApplicationModel

	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())

	// Unscoped: soft-deleted rows still occupy the unique number index.
	if err := tx.Unscoped().
		Model(&models.ArticleApplicationModel{}).
		Where("journal_id = ? AND created_at >= ? AND created_at < ?", journalID, startOfYear, endOfYear).
		Order("id DESC").
		First(&last).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	ordinal := 1
	if last.ID != 0 {
		parts := strings.Split(last.Number, "-")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid number format: %s", last.Number)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid number format: %s", last.Number)
		}
		ordinal = n + 1
	}
	return fmt.Sprintf("%03d-%03d-%02d", ordinal, journalID, year%100), nil
}
