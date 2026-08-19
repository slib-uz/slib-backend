package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ApplicationRepositoryImpl struct {
	*BaseRepository
	env *config.Config
}

// @inject
func NewApplicationRepository(baseRepository *BaseRepository, env *config.Config) repository.ApplicationRepository {
	return &ApplicationRepositoryImpl{BaseRepository: baseRepository, env: env}
}

func (this *ApplicationRepositoryImpl) Create(article *entity.ArticleCreateEntity, userId uint, technicalReviewDeadline time.Time) (*entity.ApplicationEntity, error) {

	var journal models.JournalModel
	if err := this.db().First(&journal, article.JournalID).Error; err != nil {
		return nil, _errors.Wrap(err)
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

		pendingStage := models.NewReviewStageModel(0, articleApplicationModel.ID, enum.TechnicalReviewStage, enum.StatusPending, nil, technicalReviewDeadline, nil)
		if err := tx.Create(&pendingStage).Error; err != nil {
			return err
		}

		return nil
	})

	if tr != nil {
		return nil, tr
	}

	return mapper.ApplicationModelToEntity(&articleApplicationModel), nil

}

func (this *ApplicationRepositoryImpl) generateNumber(tx *gorm.DB, journalID uint) (string, error) {
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

func (this *ApplicationRepositoryImpl) GetByJournalID(journalId uint, page, pageSize int, startDate, endDate time.Time, search string, status int) (*entity.PagingEntity[entity.ApplicationEntity], error) {
	var applications []*models.ArticleApplicationModel

	var query = this.db().Model(&models.ArticleApplicationModel{}).
		Where("article_applications.created_at BETWEEN ? AND ?", startDate, endDate)

	// If journalId is 0, don't filter by journal_id (return all applications)
	if journalId != 0 {
		query = query.Where("article_applications.journal_id = ?", journalId)
	}

	if search != "" {
		query = query.
			Joins("JOIN articles ON articles.id = article_applications.article_id").
			Where(
				"(articles.name->>'uz' ILIKE ? OR articles.name->>'ru' ILIKE ? OR articles.name->>'en' ILIKE ?)",
				"%"+search+"%", "%"+search+"%", "%"+search+"%",
			)
	}

	if status != -999 {
		query = query.
			Joins("JOIN review_stages ON review_stages.application_id = article_applications.id").
			Where("review_stages.id IN (?)", this.lastStageSubQuery()).
			Where("review_stages.status = ?", status).
			Group("article_applications.id, article_applications.created_at")
	}

	var total, _ = this.getTotalCount(query)

	if err := query.
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Preload("Article").
		Preload("User").
		Preload("Journal").
		Preload("ReviewStages", "id IN (?)", this.lastStageSubQuery()).
		Order("article_applications.created_at desc").
		Find(&applications).
		Error; err != nil {
		return nil, err
	}

	return entity.NewPagingEntity(
		page,
		pageSize,
		total,
		mapper.ApplicationModelListToEntityList(applications),
	), nil

}

func (this *ApplicationRepositoryImpl) GetByIDWithRelations(id uint) (*entity.ApplicationEntity, error) {
	var application models.ArticleApplicationModel
	if err := this.db().
		Preload("Article").
		Preload("Article.CoAuthors").
		Preload("Article.ArticleAuthorAffiliations").
		Preload("Article.ArticleAuthorAffiliations.Author").
		Preload("Article.StudyFields").
		Preload("Article.Language").
		Preload("Article.Tags").
		Preload("Journal").
		Preload("ReviewStages", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("ReviewStages.Reviewer").
		Preload("ReviewStages.AntiPlagResults").
		Preload("ReviewStages.SpellCheckResults").
		Preload("ReviewStages.AiDetectResults").
		Preload("User").
		Where("id = ?", id).
		First(&application).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := this.loadUserForAuthors(this.db(), application.Article.CoAuthors); err != nil {
		return nil, err
	}

	return mapper.ApplicationModelToEntity(&application), nil
}

func (this *ApplicationRepositoryImpl) GetUserAppByID(id, userID uint) (*entity.ApplicationEntity, error) {
	var application models.ArticleApplicationModel
	if err := this.db().
		Preload("Article").
		Preload("Article.CoAuthors").
		Preload("Article.ArticleAuthorAffiliations").
		Preload("Article.ArticleAuthorAffiliations.Author").
		Preload("Article.StudyFields").
		Preload("Article.Language").
		Preload("Article.Tags").
		Preload("Journal").
		Preload("Journal.Publisher").
		// Preload("ReviewStages.Reviewer").
		Preload("ReviewStages").
		Preload("ReviewStages.AntiPlagResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("ReviewStages.SpellCheckResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("ReviewStages.AiDetectResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&application).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := this.db().Where("article_id = ?", application.Article.ID).Find(&application.Article.References).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := this.loadUserForAuthors(this.db(), application.Article.CoAuthors); err != nil {
		return nil, err
	}
	return mapper.ApplicationModelToEntity(&application), nil
}

func (this *ApplicationRepositoryImpl) GetWithArticleAndJournal(id uint) (*entity.ApplicationEntity, error) {
	var application models.ArticleApplicationModel
	if err := this.db().
		Preload("Journal").
		Preload("Article").
		Preload("User").
		Where("id = ?", id).
		First(&application).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	if application.Article == nil {
		return nil, _errors.Wrap(gorm.ErrRecordNotFound)
	}
	return mapper.ApplicationModelToEntity(&application), nil
}

func (this *ApplicationRepositoryImpl) GetByUserID(userId uint, page, pageSize int, journalID uint) (*entity.PagingEntity[entity.ApplicationEntity], error) {
	var applications []*models.ArticleApplicationModel

	var query = this.db().Model(&models.ArticleApplicationModel{}).Where("user_id = ?", userId)

	if journalID != 0 {
		query = query.Where("journal_id = ?", journalID)
	}

	query = query.Order("created_at desc")
	var total, _ = this.getTotalCount(query)

	if err := query.
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Preload("Article").
		Preload("Journal").
		Preload("ReviewStages", "id IN (?)", this.lastStageSubQuery()).
		Find(&applications).Error; err != nil {
		return nil, err
	}

	return entity.NewPagingEntity(
		page,
		pageSize,
		total,
		mapper.ApplicationModelListToEntityList(applications),
	), nil
}

func (this *ApplicationRepositoryImpl) lastStageSubQuery() *gorm.DB {

	return this.db().
		Model(&models.ReviewStageModel{}).
		Select("MAX(review_stages.id)").
		Joins("JOIN article_applications ON article_applications.id = review_stages.application_id").
		Where("review_stages.application_id = article_applications.id AND review_stages.is_old = ?", false).
		Group("article_applications.id")
}

func (this *ApplicationRepositoryImpl) getTotalCount(query *gorm.DB) (int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, _errors.Wrap(err)
	}
	return total, nil
}

func (this *ApplicationRepositoryImpl) Update(applicationID uint, article *entity.ArticleUpdateEntity) error {
	return this.db().Transaction(func(tx *gorm.DB) error {

		var application models.ArticleApplicationModel

		if err := tx.Preload("Article").Where("id = ?", applicationID).First(&application).Error; err != nil {
			return _errors.Wrap(err)
		}

		var articleID = application.ArticleID

		oldArticle := mapper.ArticleModelToEntity(application.Article)

		var _article = mapper.ArticleUpdateEntityToModel(articleID, oldArticle, article)

		// update non zero fields
		if err := tx.Where("id = ?", articleID).Updates(_article).Error; err != nil {
			return err
		}

		// extra update
		var extraUpdate = map[string]any{}

		if article.DOI != nil {
			extraUpdate["doi"] = *article.DOI
		}
		if article.DOIClear != nil && *article.DOIClear {
			extraUpdate["doi"] = nil
		}
		if article.ExpertConclusionFile != nil {
			extraUpdate["expert_conclusion_file"] = *article.ExpertConclusionFile
		}
		if article.ExpertConclusionClear != nil && *article.ExpertConclusionClear {
			extraUpdate["expert_conclusion_file"] = nil
		}

		if err := tx.Model(_article).Where("id = ?", articleID).Updates(extraUpdate).Error; err != nil {
			return err
		}

		// Oddiy maydonlarni yangilaymiz
		if err := tx.Model(_article).
			Omit("CoAuthors", "StudyFields", "Tags").
			Updates(_article).Error; err != nil {
			return err
		}

		// Assotsiatsiyalar: faqat kelgan bo‘lsa yangilaymiz
		if article.CoAuthorsIDs != nil {
			coAuthors := make([]*models.AuthorModel, len(*article.CoAuthorsIDs))
			for i, id := range *article.CoAuthorsIDs {
				coAuthors[i] = &models.AuthorModel{Model: gorm.Model{ID: id}}
			}
			if err := tx.Model(_article).Association("CoAuthors").Replace(coAuthors); err != nil {
				return err
			}
		}

		if article.StudyFieldsIDs != nil {
			fields := make([]*models.StudyFieldModel, len(*article.StudyFieldsIDs))
			for i, id := range *article.StudyFieldsIDs {
				fields[i] = &models.StudyFieldModel{Model: gorm.Model{ID: id}}
			}
			if err := tx.Model(_article).Association("StudyFields").Replace(fields); err != nil {
				return err
			}
		}

		if article.AffiliationsIDs != nil {
			affiliations := make([]*models.ArticleAuthorAffiliationModel, len(*article.AffiliationsIDs))
			for i, affID := range *article.AffiliationsIDs {
				affiliations[i] = &models.ArticleAuthorAffiliationModel{
					Model: gorm.Model{ID: affID},
				}
			}
			if err := tx.Model(_article).Association("ArticleAuthorAffiliations").Replace(affiliations); err != nil {
				return err
			}
		}

		if article.Tags != nil {
			tagIDs, err := getOrCreateTagsByLang(tx, article.Tags)
			if err != nil {
				return err
			}
			tags := make([]*models.TagModel, len(tagIDs))
			for i, tagID := range tagIDs {
				tags[i] = &models.TagModel{ID: tagID}
			}
			if err := tx.Model(_article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		var currentStage models.ReviewStageModel
		if err := tx.Where("application_id = ? AND stage = ? AND status = ? AND is_old = ?", applicationID, enum.TechnicalReviewStage, enum.StatusRejected, false).
			Order("created_at DESC").
			First(&currentStage).
			Error; err != nil {
			return _errors.Wrap(err)
		}

		if err := tx.Model(&currentStage).Update("is_old", true).Error; err != nil {
			return err
		}

		// new stage
		newStage := models.NewReviewStageModel(0, applicationID, enum.TechnicalReviewStage, enum.StatusPending, &currentStage.ID, utils.MakeDeadline(time.Now(), this.env.ReviewDeadlineDays), nil)

		if err := tx.Create(&newStage).Error; err != nil {
			return _errors.Wrap(err)
		}

		return nil
	})
}

func (this *ApplicationRepositoryImpl) FindByID(id uint) (*entity.ApplicationEntity, error) {
	var application models.ArticleApplicationModel
	if err := this.db().
		Where("id = ?", id).
		Preload("Article").
		Preload("Article.Journal").
		Preload("Article.Language").
		Preload("Article.CoAuthors").
		Preload("Article.StudyFields").
		First(&application).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ApplicationModelToEntity(&application), nil
}

func (this *ApplicationRepositoryImpl) GetCheckForPaymentByApplicationID(applicationID uint) (*entity.ReviewStageEntity, error) {
	var stage models.ReviewStageModel

	if err := this.db().
		Preload("Application").
		Preload("Application.Journal", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "publishing_price", "name", "short_name", "issn_paper", "issn_online")
		}).
		Where("application_id = ? AND id IN (?)", applicationID, this.lastStageSubQuery()).
		First(&stage).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.ReviewStageModelToEntity(&stage), nil
}

func (this *ApplicationRepositoryImpl) Publish(articleID, applicationID uint, finalFilePath string) error {
	var article = models.ArticleModel{
		Model:           gorm.Model{ID: articleID, UpdatedAt: time.Now()},
		ContentFilePath: finalFilePath,
		PublicationDate: datatypes.Date(time.Now()),
		IsPublished:     true,
	}

	return this.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ArticleModel{}).Where("id = ?", articleID).Updates(article).Error; err != nil {
			return _errors.Wrap(err)
		}
		err := tx.Model(&models.ArticleApplicationModel{}).Where("id = ?", applicationID).Update("is_published", true).Error
		return _errors.Wrap(err)
	})

}

func (this *ApplicationRepositoryImpl) UpdateFile(applicationID uint, filePath string) error {
	var application models.ArticleApplicationModel

	if err := this.db().
		Preload("Article").
		Where("id = ?", applicationID).
		First(&application).Error; err != nil {
		return _errors.Wrap(err)
	}

	return this.db().
		Model(&application.Article).
		Update("content_file_path", filePath).
		Error
}

func (this *ApplicationRepositoryImpl) FindByIDWithRelations(id uint) (*entity.ApplicationEntity, error) {
	var application models.ArticleApplicationModel
	if err := this.db().
		Preload("Article").
		Preload("Article.Journal").
		Preload("Article.Language").
		Preload("Article.CoAuthors").
		Preload("Article.StudyFields").
		Preload("Article.Tags").
		Preload("User").
		Where("id = ?", id).
		First(&application).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	for _, field := range application.Article.StudyFields {
		if err := this.LoadParentsRecursive(this.db(), field); err != nil {
			return nil, _errors.Wrap(err)
		}
	}

	return mapper.ApplicationModelToEntity(&application), nil
}

func (this *ApplicationRepositoryImpl) LoadParentsRecursive(db *gorm.DB, field *models.StudyFieldModel) error {
	if field.ParentID == nil {
		return nil
	}

	var parent models.StudyFieldModel
	if err := db.Preload("Parent").First(&parent, *field.ParentID).Error; err != nil {
		return err
	}

	field.Parent = &parent
	return this.LoadParentsRecursive(db, &parent)
}

func (this *ApplicationRepositoryImpl) CheckUniqueConstraints(article *entity.ArticleCreateEntity) error {
	var count int64

	if err := this.db().Model(&models.ArticleModel{}).
		Where("doi = ?", article.DOI).
		Count(&count).Error; err != nil {
		return _errors.Wrap(err)
	}

	if count > 0 {
		return errors.New("article with this DOI or ROI already exists")
	}

	return nil
}

// func (this *ApplicationRepositoryImpl) chunkStrings(ss []string, size int) [][]string {
// 	if size <= 0 {
// 		size = 1000
// 	}
// 	var chunks [][]string
// 	for i := 0; i < len(ss); i += size {
// 		end := i + size
// 		if end > len(ss) {
// 			end = len(ss)
// 		}
// 		chunks = append(chunks, ss[i:end])
// 	}
// 	return chunks
// }

func (this *ApplicationRepositoryImpl) loadUserForAuthors(db *gorm.DB, authors []*models.AuthorModel) error {

	var ids []string
	for _, author := range authors {
		if author == nil {
			continue
		}
		ids = append(ids, author.ScienceID)
	}

	if len(ids) == 0 {
		return nil
	}

	users, err := this.findUsersByIDs(db, ids)
	if err != nil {
		return err
	}

	this.linkUsersToAuthors(authors, users)

	return nil
}

func (this *ApplicationRepositoryImpl) findUsersByIDs(db *gorm.DB, ids []string) ([]models.UserModel, error) {
	var users []models.UserModel
	err := db.Where("science_id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, _errors.Wrap(err)
	}
	return users, nil
}

func (this *ApplicationRepositoryImpl) linkUsersToAuthors(authors []*models.AuthorModel, users []models.UserModel) {
	userMap := make(map[string]models.UserModel, len(users))

	for _, user := range users {
		if _, exists := userMap[user.ScienceID]; exists {
			continue
		}
		userMap[user.ScienceID] = user
	}

	for _, author := range authors {
		if author == nil {
			continue
		}

		user, found := userMap[author.ScienceID]

		if !found {
			author.Users = nil
			continue
		}

		usersList := []models.UserModel{user}
		author.Users = &usersList
	}
}

func (this *ApplicationRepositoryImpl) GetApplicationCountByUserID(userID uint, isPublished bool) (int64, error) {
	var count int64
	if err := this.db().Model(&models.ArticleApplicationModel{}).Where("user_id = ? AND is_published = ?", userID, isPublished).Count(&count).Error; err != nil {
		return 0, _errors.Wrap(err)
	}
	return count, nil
}
