package repository

import (
	"errors"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ArticleRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewArticleRepository(repository *BaseRepository) repository.ArticleRepository {
	return &ArticleRepositoryImpl{BaseRepository: repository}
}

func (this *ArticleRepositoryImpl) GetAuthorsCount() (int64, error) {
	var count int64
	err := this.db().Model(&models.ArticleModel{}).
		Where("articles.is_published = ?", true).
		Joins(`
			JOIN article_co_authors
			  ON article_co_authors.article_model_id = articles.id
		`).
		Distinct("article_co_authors.author_model_id").
		Count(&count).Error

	return count, err
}

func (this *ArticleRepositoryImpl) GetArticleCount() (int64, error) {
	var count int64
	if err := this.database.GormDB.Model(&models.ArticleModel{}).Where("is_published = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (this *ArticleRepositoryImpl) GetArticleCountByUserID(userID uint, isPublished bool) (int64, error) {
	var count int64
	err := this.db().
		Model(&models.ArticleModel{}).
		Joins("JOIN article_co_authors ON article_co_authors.article_model_id = articles.id").
		Joins("JOIN authors ON authors.id = article_co_authors.author_model_id").
		Joins("JOIN users ON users.science_id = authors.science_id").
		Where("users.id = ? AND users.deleted_at IS NULL", userID).
		Where("articles.is_published = ?", isPublished).
		Where("articles.deleted_at IS NULL").
		Distinct("articles.id").
		Count(&count).Error

	if err != nil {
		return 0, _errors.Wrap(err)
	}
	return count, nil
}

func (this *ArticleRepositoryImpl) GetArticleCountByEditionID(editionID uint) (int64, error) {
	var count int64
	err := this.db().
		Model(&models.ArticleModel{}).
		Where("edition_id = ?", editionID).
		Count(&count).Error

	if err != nil {
		return 0, _errors.Wrap(err)
	}
	return count, nil
}

func (this *ArticleRepositoryImpl) GetByIDWithJournal(id uint) (*entity.ArticleBasicEntity, error) {
	var article models.ArticleModel
	if err := this.db().
		Preload("Journal").
		// Preload("Tags").
		Where("id = ?", id).
		First(&article).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ArticleModelToBasicEntity(&article), nil
}

func (this *ArticleRepositoryImpl) FindByID(id uint) (*entity.ArticleEntity, error) {
	var article models.ArticleModel
	if err := this.db().
		Where("id = ?", id).
		First(&article).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ArticleModelToEntity(&article), nil
}

func (this *ArticleRepositoryImpl) FindByIDWithAuthors(id uint) (*entity.ArticleEntity, error) {
	var article models.ArticleModel
	if err := this.db().
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Where("id = ?", id).
		First(&article).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ArticleModelToEntity(&article), nil
}

func (this *ArticleRepositoryImpl) FindByIDsWithAuthors(ids []uint) ([]*entity.ArticleEntity, error) {
	var articles []*models.ArticleModel
	if err := this.db().
		Preload("CoAuthors").
		Where("id IN ?", ids).
		Find(&articles).
		Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.ArticleModelListToEntityList(articles), nil
}

func (this *ArticleRepositoryImpl) UpdateViewsCount(counts map[uint]int64) error {

	if len(counts) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(counts))     // WHERE id IN (…)
	params := make([]any, 0, len(counts)*2) // CASE parametrlari
	caseSQL := "CASE id"

	for id, d := range counts {
		ids = append(ids, id)

		//  WHEN <id> THEN views_count + <delta>
		caseSQL += " WHEN ? THEN views_count + ?::bigint"
		params = append(params, id, d)
	}
	caseSQL += " ELSE views_count END" // boshqa id‑lar o‘zgarmaydi

	return this.db().
		Model(&models.ArticleModel{}).
		Where("id IN ?", ids).
		Update("views_count", gorm.Expr(caseSQL, params...)).
		Error
}

func (this *ArticleRepositoryImpl) UpdateROI(
	id uint,
	externalID uint,
	roi *string,
) error {

	tx := this.db().
		Model(&models.ArticleModel{}).
		Where("id = ?", id).
		Where(`
			NOT EXISTS (
				SELECT 1
				FROM articles a2
				WHERE a2.external_id = ?
				  AND a2.id <> ?
				  AND a2.deleted_at IS NULL
			)
		`, externalID, id).
		Updates(map[string]any{
			"roi":         roi,
			"external_id": externalID,
		})

	if tx.Error != nil {
		return _errors.Wrap(tx.Error)
	}

	if tx.RowsAffected == 0 {
		return _errors.Wrap(fmt.Errorf("external_id already in use"))
	}

	return nil
}

// PublishedIDsWithoutROI — ROI olinmagan (roi IS NULL) published maqolalarning ID'larini qaytaradi.
func (this *ArticleRepositoryImpl) PublishedIDsWithoutROI() ([]uint, error) {
	var ids []uint

	if err := this.db().
		Model(&models.ArticleModel{}).
		Where("is_published = ? AND roi IS NULL", true).
		Order("id ASC").
		Pluck("id", &ids).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return ids, nil
}

func (this *ArticleRepositoryImpl) UpdateDOI(id uint, doi string) error {
	if err := this.db().
		Model(&models.ArticleModel{}).
		Where("id = ?", id).
		Update("doi", doi).Error; err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *ArticleRepositoryImpl) UpdateStatus(id uint, isSent bool) error {
	if err := this.db().
		Model(&models.ArticleModel{}).
		Where("id = ?", id).
		Update("is_sent", isSent).Error; err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *ArticleRepositoryImpl) ManageUpdate(id uint, e *entity.ArticleManageUpdateEntity) error {
	return this.db().Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"name":                   mapper.ToGormJson(e.Name),
			"co_authors_count":       e.CoAuthorsCount,
			"language_id":            e.LanguageID,
			"content_file_path":      e.ContentFile,
			"annotation":             mapper.ToGormJson(e.Annotation),
			"access_type":            e.AccessType,
			"doi":                    e.DOI,
			"expert_conclusion_file": e.ExpertConclusionFile,
			"pages":                  e.Pages,
			"web_address":            e.WebAddress,
			"publication_date":       datatypes.Date(*e.PublishedDate),
			"unconfirmed_authors":    mapper.ToGormJson(e.UnconfirmedAuthors),
		}

		if err := tx.Model(&models.ArticleModel{}).
			Where("id = ?", id).
			Updates(updates).Error; err != nil {
			return _errors.Wrap(err)
		}

		var article models.ArticleModel
		if err := tx.First(&article, id).Error; err != nil {
			return _errors.Wrap(err)
		}

		coAuthors := make([]*models.AuthorModel, len(e.CoAuthorsIDs))
		for i, authorID := range e.CoAuthorsIDs {
			coAuthors[i] = &models.AuthorModel{Model: gorm.Model{ID: authorID}}
		}
		if err := tx.Model(&article).Association("CoAuthors").Replace(coAuthors); err != nil {
			return _errors.Wrap(err)
		}

		studyFields := make([]*models.StudyFieldModel, len(e.StudyFieldsIDs))
		for i, fieldID := range e.StudyFieldsIDs {
			studyFields[i] = &models.StudyFieldModel{Model: gorm.Model{ID: fieldID}}
		}
		if err := tx.Model(&article).Association("StudyFields").Replace(studyFields); err != nil {
			return _errors.Wrap(err)
		}

		tagIDs, err := getOrCreateTagsByLang(tx, e.Tags)
		if err != nil {
			return err
		}
		tags := make([]*models.TagModel, len(tagIDs))
		for i, tagID := range tagIDs {
			tags[i] = &models.TagModel{ID: tagID}
		}
		if err := tx.Model(&article).Association("Tags").Replace(tags); err != nil {
			return _errors.Wrap(err)
		}

		affiliations := make([]*models.ArticleAuthorAffiliationModel, len(e.AffiliationsIDs))
		for i, affID := range e.AffiliationsIDs {
			affiliations[i] = &models.ArticleAuthorAffiliationModel{
				Model: gorm.Model{ID: affID},
			}
		}
		if err := tx.Model(&article).Association("ArticleAuthorAffiliations").Replace(affiliations); err != nil {
			return _errors.Wrap(err)
		}

		if e.References != nil {
			if err := tx.Where("article_id = ?", id).Delete(&models.ReferenceModel{}).Error; err != nil {
				return _errors.Wrap(err)
			}

			if len(e.References) > 0 {
				var references []*models.ReferenceModel
				for _, refName := range e.References {
					references = append(references, models.NewReferenceModel(refName, id))
				}
				if err := tx.Create(&references).Error; err != nil {
					return _errors.Wrap(err)
				}
			}
		}

		return nil
	})
}

func (this *ArticleRepositoryImpl) AddCoAuthor(articleID uint, authorID uint) error {
	return this.db().Transaction(func(tx *gorm.DB) error {
		var article models.ArticleModel
		if err := tx.First(&article, articleID).Error; err != nil {
			return _errors.Wrap(err)
		}

		author := models.AuthorModel{Model: gorm.Model{ID: authorID}}

		if err := tx.Model(&article).Association("CoAuthors").Append(&author); err != nil {
			return _errors.Wrap(err)
		}

		count := tx.Model(&article).Association("CoAuthors").Count()
		if err := tx.Model(&article).Update("co_authors_count", count).Error; err != nil {
			return _errors.Wrap(err)
		}

		return nil
	})
}

func (this *ArticleRepositoryImpl) OnlyArticleNameByLegacyAuthorIDs(ids []uint, page, pageSize int) (*entity.PagingEntity[entity.ArticleOnlyNameEntity], error) {
	var articles []*models.ArticleModel
	var total int64

	query := this.db().
		Model(&models.ArticleModel{}).
		Joins("JOIN legacy_author_articles ON legacy_author_articles.article_model_id = articles.id").
		Where("legacy_author_articles.legacy_author_id IN ?", ids)

	if err := query.Distinct("articles.id").Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	err := query.
		Select("articles.id, articles.name").
		Distinct("articles.id, articles.name").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&articles).Error

	if err != nil {
		return nil, _errors.Wrap(err)
	}

	return entity.NewPagingEntity(page, pageSize, total, mapper.ArticleModelListToOnlyNameEntityList(articles)), nil
}

func (this *ArticleRepositoryImpl) CreatePublishArticle(entity *entity.ArticlePublishEntity, journalID uint) (uint, error) {
	var createdArticleID uint

	err := this.db().Transaction(func(tx *gorm.DB) error {
		article := mapper.ArticlePublishEntityToArticleModel(entity, journalID)
		article.IsPublished = true
		article.PublicationDate = datatypes.Date(*entity.PublishedDate)

		if err := tx.Create(article).Error; err != nil {
			return _errors.Wrap(err)
		}

		createdArticleID = article.ID

		if err := this.attachStudyFields(tx, article, entity.StudyFieldsIDs); err != nil {
			return err
		}

		if err := this.attachTags(tx, article, entity.Tags); err != nil {
			return err
		}

		if err := this.attachCoAuthors(tx, article, entity.CoAuthorsIDs); err != nil {
			return err
		}

		if err := this.attachAffiliations(tx, article.ID, entity.ArticleAuthorAffiliationIDs); err != nil {
			return err
		}

		return this.createReferences(tx, article.ID, entity.References)
	})

	return createdArticleID, err
}

func (this *ArticleRepositoryImpl) CreatePublishArticleWithAffiliations(entity *entity.ArticlePublishEntity, journalID uint) (uint, error) {
	var createdArticleID uint

	err := this.db().Transaction(func(tx *gorm.DB) error {
		article := mapper.ArticlePublishEntityToArticleModel(entity, journalID)
		article.IsPublished = true
		article.PublicationDate = datatypes.Date(*entity.PublishedDate)

		if err := tx.Create(article).Error; err != nil {
			return _errors.Wrap(err)
		}

		createdArticleID = article.ID

		if err := this.attachStudyFields(tx, article, entity.StudyFieldsIDs); err != nil {
			return err
		}

		if err := this.attachTags(tx, article, entity.Tags); err != nil {
			return err
		}

		if err := this.attachCoAuthors(tx, article, coAuthorIDsFromAffiliations(entity.AuthorAffiliations)); err != nil {
			return err
		}

		if err := this.createAffiliationsWithOrgs(tx, article.ID, entity.AuthorAffiliations); err != nil {
			return err
		}

		return this.createReferences(tx, article.ID, entity.References)
	})

	return createdArticleID, err
}

func (this *ArticleRepositoryImpl) createAffiliationsWithOrgs(tx *gorm.DB, articleID uint, affiliations []*entity.ArticleAuthorAffiliationEntity) error {
	if len(affiliations) == 0 {
		return nil
	}

	affiliationModels := make([]*models.ArticleAuthorAffiliationModel, 0, len(affiliations))
	for _, aff := range affiliations {
		tin := utils.NormalizeTin(aff.OrganizationTin)
		var org models.OrganizationModel
		err := tx.
			Where(models.OrganizationModel{Tin: tin}).
			Attrs(models.OrganizationModel{Name: aff.OrganizationName}).
			FirstOrCreate(&org).Error
		if err != nil {
			wrapped := _errors.Wrap(err)
			if !errors.Is(wrapped, response.ConflictError) {
				return wrapped
			}
			if lookupErr := tx.Where("tin = ?", tin).First(&org).Error; lookupErr != nil {
				return _errors.Wrap(lookupErr)
			}
		}

		organizationID := org.ID
		affiliationModels = append(affiliationModels, &models.ArticleAuthorAffiliationModel{
			ArticleID:        &articleID,
			AuthorID:         aff.AuthorID,
			OrganizationID:   &organizationID,
			OrganizationName: org.Name,
			OrganizationTin:  org.Tin,
			PositionName:     aff.PositionName,
		})
	}

	return _errors.Wrap(tx.Create(&affiliationModels).Error)
}

// coAuthorIDsFromAffiliations returns the unique author IDs referenced by the
// affiliations, preserving their order. Used to keep the article's co-authors in
// sync with the inline affiliations for the integration flow.
func coAuthorIDsFromAffiliations(affiliations []*entity.ArticleAuthorAffiliationEntity) []uint {
	seen := make(map[uint]bool, len(affiliations))
	ids := make([]uint, 0, len(affiliations))
	for _, aff := range affiliations {
		if aff.AuthorID == 0 || seen[aff.AuthorID] {
			continue
		}
		seen[aff.AuthorID] = true
		ids = append(ids, aff.AuthorID)
	}
	return ids
}

func (this *ArticleRepositoryImpl) attachCoAuthors(tx *gorm.DB, article *models.ArticleModel, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	authors := make([]*models.AuthorModel, 0, len(ids))
	for _, id := range ids {
		authors = append(authors, &models.AuthorModel{Model: gorm.Model{ID: id}})
	}
	return _errors.Wrap(tx.Model(article).Association("CoAuthors").Append(authors))
}

func (this *ArticleRepositoryImpl) attachAffiliations(tx *gorm.DB, articleID uint, affIDs []uint) error {
	if len(affIDs) == 0 {
		return nil
	}

	err := tx.Model(&models.ArticleAuthorAffiliationModel{}).
		Where("id IN ?", affIDs).
		UpdateColumn("article_id", articleID).Error

	if err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *ArticleRepositoryImpl) attachStudyFields(tx *gorm.DB, article *models.ArticleModel, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	fields := make([]*models.StudyFieldModel, 0, len(ids))
	for _, id := range ids {
		m := models.NewStudyFieldModel(id, nil, nil, nil)
		fields = append(fields, &m)
	}

	err := tx.Model(article).Association("StudyFields").Append(fields)
	return _errors.Wrap(err)
}

func (this *ArticleRepositoryImpl) attachTags(tx *gorm.DB, article *models.ArticleModel, tags entity.TagNamesByLang) error {
	tagIDs, err := getOrCreateTagsByLang(tx, tags)
	if err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}

	modelsToAttach := make([]*models.TagModel, len(tagIDs))
	for i, tagID := range tagIDs {
		modelsToAttach[i] = &models.TagModel{ID: tagID}
	}

	return _errors.Wrap(tx.Model(article).Association("Tags").Append(modelsToAttach))
}

func (this *ArticleRepositoryImpl) createReferences(tx *gorm.DB, articleID uint, refs []string) error {
	if len(refs) == 0 {
		return nil
	}

	modelsRefs := make([]*models.ReferenceModel, 0, len(refs))
	for _, name := range refs {
		modelsRefs = append(modelsRefs, models.NewReferenceModel(name, articleID))
	}
	return _errors.Wrap(tx.Create(&modelsRefs).Error)
}

func (this *ArticleRepositoryImpl) CheckDOIExists(doi string) (bool, error) {
	var count int64
	err := this.db().Model(&models.ArticleModel{}).Where("doi = ?", doi).Count(&count).Error
	return count > 0, _errors.Wrap(err)
}

func (this *ArticleRepositoryImpl) CheckROIExists(roi string) (bool, error) {
	var count int64
	err := this.db().Model(&models.ArticleModel{}).Where("roi = ?", roi).Count(&count).Error
	return count > 0, _errors.Wrap(err)
}

func (this *ArticleRepositoryImpl) CheckExternalIDExists(externalID uint) (bool, error) {
	var count int64
	err := this.db().Model(&models.ArticleModel{}).Where("external_id = ?", externalID).Count(&count).Error
	return count > 0, _errors.Wrap(err)
}
