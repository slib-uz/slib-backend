package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

type PublishedArticleRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewPublishedArticleRepository(repository *BaseRepository) repository.PublishedArticleRepository {
	return &PublishedArticleRepositoryImpl{BaseRepository: repository}
}

func (this *PublishedArticleRepositoryImpl) model() *gorm.DB {
	return this.BaseRepository.db().Where("is_published = ?", true)
}

func (this *PublishedArticleRepositoryImpl) GetByIDWithRelations(id uint) (*entity2.ArticleEntity, error) {

	var article models.ArticleModel

	if err := this.model().
		Preload("Tags").
		Preload("Journal").
		Preload("Journal.Publisher").
		Preload("Language").
		Preload("CoAuthors").
		Preload("StudyFields").
		Preload("Edition").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Where("id = ? AND is_published = ?", id, true).
		First(&article).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	// Load References separately
	var references []*models.ReferenceModel
	if err := this.db().Where("article_id = ?", id).Find(&references).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	article.References = references

	return mapper.ArticleModelToEntity(&article), nil

}

func (this *PublishedArticleRepositoryImpl) GetByAuthor(scienceID string, page, pageSize int, search string, studyFieldIds []uint, year int, journalID *uint) (*entity2.PagingEntity[entity2.ArticleEntity], error) {

	var author models.AuthorModel
	err := this.db().
		Where("science_id = ?", scienceID).
		First(&author).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity2.NewPagingEntity(
				page,
				pageSize,
				0,
				[]*entity2.ArticleEntity{},
			), nil
		}
		return nil, infraError.Wrap(err)
	}

	var articles []*models.ArticleModel

	query := this.model().
		Omit("Annotation").
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Preload("StudyFields").
		Preload("Tags").
		Preload("Journal").
		Preload("Language").
		Where("is_published = ?", true)

	query = query.
		Joins("JOIN article_co_authors ON article_co_authors.article_model_id = articles.id").
		Where("article_co_authors.author_model_id = ?", author.ID)

	if journalID != nil {
		query = query.Where("articles.journal_id = ?", *journalID)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"(name->>'uz' ILIKE ? OR name->>'ru' ILIKE ? OR name->>'en' ILIKE ?)",
			like, like, like,
		)
	}

	if len(studyFieldIds) > 0 {
		query = query.
			Joins(`JOIN article_study_fields asf ON asf.article_model_id = articles.id`).
			Where("asf.study_field_model_id IN ?", studyFieldIds).
			Distinct("articles.id")
	}

	if year > 0 {
		query = query.Where("EXTRACT(YEAR FROM publication_date) = ?", year)
	}

	var total int64
	if err := query.Model(&models.ArticleModel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).
		Error; err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(
		page,
		pageSize,
		total,
		mapper.ArticleModelListToEntityList(articles),
	), nil
}

func (this *PublishedArticleRepositoryImpl) GetByJournal(journalID uint, page, pageSize int, search string, authorSearch, tag string, studyFieldIds []uint, year int) (*entity2.PagingEntity[entity2.ArticleEntity], error) {
	var articles []*models.ArticleModel

	query := this.model().Model(&models.ArticleModel{}).Where("journal_id = ?", journalID)
	if search != "" {
		query = query.Where(
			"(name->>'uz' ILIKE ? OR name->>'ru' ILIKE ? OR name->>'en' ILIKE ?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if authorSearch != "" {
		like := "%" + authorSearch + "%"
		query = query.
			Joins(`JOIN article_co_authors aca ON aca.article_model_id = articles.id`).
			Joins(`JOIN authors a ON a.id = aca.author_model_id`).
			Where("(a.full_name ILIKE ? OR a.science_id ILIKE ?)", like, like).
			Distinct("articles.id")
	}

	if tag != "" {
		like := "%" + tag + "%"
		query = query.
			Joins(`JOIN article_tags at ON at.article_model_id = articles.id`).
			Joins(`JOIN tags t ON t.id = at.tag_model_id`).
			Where("t.name ILIKE ?", like).
			Distinct("articles.id")
	}

	if len(studyFieldIds) > 0 {
		query = query.
			Joins(`Join article_study_fields asf ON asf.article_model_id = articles.id`).
			Where("asf.study_field_model_id IN ?", studyFieldIds).
			Distinct("articles.id")
	}

	if year > 0 {
		query = query.Where("EXTRACT(YEAR FROM publication_date) = ?", year)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Omit("Annotation").
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Preload("StudyFields").
		Preload("Tags").
		Preload("Journal").
		Preload("Language").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).
		Error; err != nil {
		return nil, err
	}
	return mapper.PagingMapper(
		articles,
		mapper.ArticleModelToEntity,
		page,
		pageSize,
		total,
	), nil
}

func (this *PublishedArticleRepositoryImpl) GetTopArticles(page, pageSize int) (*entity2.PagingEntity[entity2.ArticleEntity], error) {
	var articles []*models.ArticleModel
	var total int64

	if err := this.db().
		Model(&models.ArticleModel{}).
		Where("is_published = ?", true).
		Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := this.db().
		Model(&models.ArticleModel{}).
		Where("articles.is_published = ?", true).
		Preload("Journal").
		Preload("Language").
		Preload("StudyFields").
		Preload("Tags").
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Order("views_count DESC, rating_sum DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return entity2.NewPagingEntity(
		page,
		pageSize,
		total,
		mapper.ArticleModelListToEntityList(articles),
	), nil
}

func (this *PublishedArticleRepositoryImpl) GetByEdition(journalID, editionID uint, unassigned bool, page, pageSize int) (*entity2.PagingEntity[entity2.ArticleEntity], error) {
	var articles []*models.ArticleModel
	var total int64

	query := this.model()

	if journalID > 0 {
		query = query.Where("journal_id = ?", journalID)
	}

	if unassigned {
		query = query.Where("edition_id IS NULL")
	} else if editionID > 0 {
		query = query.Where("edition_id = ?", editionID)
	}

	if err := query.Model(&models.ArticleModel{}).Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.
		Omit("Annotation").
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author").
		Preload("StudyFields").
		Preload("Tags").
		Preload("Journal").
		Preload("Language").
		Order("publication_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).
		Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(articles, mapper.ArticleModelToEntity, page, pageSize, total), nil
}

// joinOr faqat shu fayldagi konstanta SQL bo'laklarini birlashtiradi.
// Foydalanuvchi kiritishi bu yerga tushmaydi — u args orqali ? ga bog'lanadi.
//
// Alohida funksiyaga chiqarilgani bejiz emas: chaqiruv joyida bu ifoda
// xavfli koddan tashqi ko'rinishi bilan farq qilmasdi.
func joinOr(conditions []string) string {
	return "(" + strings.Join(conditions, " OR ") + ")"
}

// ArticleSortFields — /api/articles/list uchun ruxsat etilgan tartiblash
// maydonlari. Nomlar articles_list_handler.go dagi Enums(...) ro'yxatiga mos.
//
// Ustunlar jadval nomi bilan to'liq yozilgan: qidiruv parametrlari berilganda
// so'rov journals jadvalini JOIN qiladi, va views_count/rating_sum ikkala
// jadvalda ham bor.
var ArticleSortFields = sorting.New("articles.publication_date DESC", map[string]string{
	"views_count":      "articles.views_count",
	"rating_sum":       "articles.rating_sum",
	"publication_date": "articles.publication_date",
})

func (this *PublishedArticleRepositoryImpl) GetAll(
	fromYear, toYear *int, journal uint, languages []uint,
	studyfields []uint, tags []string,
	articleTitle, authorName, journalName, issnPaper, issnOnline *string,
	sort string,
	page, pageSize int,
	accessType enum.AccessType,
	hasDoi *bool,
) (*entity2.PagingEntity[entity2.ArticleEntity], error) {
	var articles []*models.ArticleModel

	// Tartiblash birinchi tekshiriladi: yaroqsiz qiymat uchun bazaga
	// murojaat qilishning hojati yo'q.
	orderExpr, err := ArticleSortFields.Resolve(sort)
	if err != nil {
		return nil, err
	}

	query := this.db().Model(&models.ArticleModel{}).
		Where("articles.is_published = ?", true).
		Preload("Journal").
		//Preload("Language").
		//Preload("StudyFields").
		//Preload("Tags").
		Preload("CoAuthors").
		Preload("ArticleAuthorAffiliations").
		Preload("ArticleAuthorAffiliations.Author")

	if fromYear != nil {
		query = query.Where("EXTRACT(YEAR FROM publication_date) >= ?", *fromYear)
	}
	if toYear != nil {
		query = query.Where("EXTRACT(YEAR FROM publication_date) <= ?", *toYear)
	}
	if journal != 0 {
		query = query.Where("articles.journal_id = ?", journal)
	}
	if len(languages) > 0 {
		query = query.Where("articles.language_id IN ?", languages)
	}
	if accessType != 0 {
		query = query.Where("articles.access_type = ?", accessType)
	}
	if hasDoi != nil {
		if *hasDoi {
			query = query.Where("articles.doi IS NOT NULL AND articles.doi != ''")
		} else {
			query = query.Where("articles.doi IS NULL OR articles.doi = ''")
		}
	}

	if len(studyfields) > 0 {
		query = query.
			Joins("JOIN article_study_fields ON article_study_fields.article_model_id = articles.id").
			Where("article_study_fields.study_field_model_id IN ?", studyfields)
	}

	if len(tags) > 0 {
		query = query.
			Joins(`
                INNER JOIN article_tags ON article_tags.article_model_id = articles.id
                INNER JOIN tags ON tags.id = article_tags.tag_model_id
            `)

		for i, tag := range tags {
			if i == 0 {
				query = query.Where("tags.name ILIKE ?", "%"+tag+"%")
			} else {
				query = query.Or("tags.name ILIKE ?", "%"+tag+"%")
			}
		}
	}

	hasSearch := articleTitle != nil || authorName != nil || journalName != nil || issnPaper != nil || issnOnline != nil
	if hasSearch {
		query = query.Joins(`
            LEFT JOIN article_co_authors
              ON article_co_authors.article_model_id = articles.id
            LEFT JOIN authors
              ON authors.id = article_co_authors.author_model_id
            LEFT JOIN journals j
              ON j.id = articles.journal_id
        `)

		var conditions []string
		var args []interface{}

		if articleTitle != nil {
			like := "%" + *articleTitle + "%"
			conditions = append(conditions, "(articles.name->>'uz' ILIKE ? OR articles.name->>'ru' ILIKE ? OR articles.name->>'en' ILIKE ?)")
			args = append(args, like, like, like)
		}
		if authorName != nil {
			like := "%" + *authorName + "%"
			conditions = append(conditions, "authors.full_name ILIKE ?")
			args = append(args, like)
		}
		if journalName != nil {
			like := "%" + *journalName + "%"
			conditions = append(conditions, "(j.name->>'uz' ILIKE ? OR j.name->>'ru' ILIKE ? OR j.name->>'en' ILIKE ?)")
			args = append(args, like, like, like)
		}
		if issnPaper != nil {
			like := "%" + *issnPaper + "%"
			conditions = append(conditions, "j.issn_paper ILIKE ?")
			args = append(args, like)
		}
		if issnOnline != nil {
			like := "%" + *issnOnline + "%"
			conditions = append(conditions, "j.issn_online ILIKE ?")
			args = append(args, like)
		}

		if len(conditions) > 0 {
			query = query.Where(joinOr(conditions), args...)
		}
	}

	var total int64
	if err := query.Session(&gorm.Session{}).
		Distinct("articles.id").
		Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	query = query.Order(orderExpr)

	query = query.Select("articles.*").Group("articles.id")

	if err := query.
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&articles).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	return mapper.PagingMapper(articles, mapper.ArticleModelToEntity, page, pageSize, total), nil
}
