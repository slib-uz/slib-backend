package repository

import (
	"errors"
	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraErrors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AuthorRepositoryImpl struct {
	database           *db.Database
	researchMetricRepo repository.ResearchMetricRepository
}

// @inject
func NewAuthorRepository(database *db.Database, researchMetricRepo repository.ResearchMetricRepository) repository.AuthorRepository {
	return &AuthorRepositoryImpl{
		database:           database,
		researchMetricRepo: researchMetricRepo,
	}
}

func (r *AuthorRepositoryImpl) GetTopAuthorsWithArticleCount(top int, journalID *uint) ([]*entity.AuthorEntity, error) {
	type AuthorWithArticleCount struct {
		models.AuthorModel
		ArticleCount int64
		UserID       *uint
		Photo        *string
	}

	if top <= 0 {
		top = 10
	}

	var authors []AuthorWithArticleCount

	query := r.db().
		Table("authors a").
		Select("a.*, COUNT(DISTINCT ac.article_model_id) as article_count, u.id as user_id, u.photo").
		Joins("JOIN article_co_authors ac ON a.id = ac.author_model_id").
		Joins("JOIN articles ar ON ac.article_model_id = ar.id AND ar.is_published = TRUE").
		Joins("LEFT JOIN users u ON u.science_id = a.science_id AND u.deleted_at IS NULL").
		Where("a.deleted_at IS NULL")

	if journalID != nil {
		query = query.Where("ar.journal_id = ?", *journalID)
	}

	err := query.
		Group("a.id, u.id, u.photo").
		Order("article_count DESC").
		Limit(top).
		Scan(&authors).Error

	if err != nil {
		return nil, err
	}

	if len(authors) == 0 {
		return []*entity.AuthorEntity{}, nil
	}

	userIDs := make([]uint, 0, len(authors))
	for _, author := range authors {
		if author.UserID != nil {
			userIDs = append(userIDs, *author.UserID)
		}
	}

	userMap := make(map[uint]*models.UserModel)
	if len(userIDs) > 0 {
		var users []models.UserModel
		if err := r.db().Where("id IN ? AND deleted_at IS NULL", userIDs).Find(&users).Error; err == nil {
			for i := range users {
				userMap[users[i].ID] = &users[i]
			}
		}
	}

	jobMap := make(map[uint]*models.JobModel)
	if len(userIDs) > 0 {
		var jobs []models.JobModel
		sub := r.db().Table("jobs").
			Select("DISTINCT ON (user_id) *").
			Where("user_id IN ? AND deleted_at IS NULL", userIDs).
			Order("user_id, created_at DESC")
		_ = r.db().Table("(?) as j", sub).Find(&jobs).Error

		for i := range jobs {
			jobMap[jobs[i].UserID] = &jobs[i]
		}
	}

	researchMetricsMap := make(map[uint][]*entity.ResearchMetricEntity)
	if len(userIDs) > 0 {
		metrics, err := r.researchMetricRepo.GetByUserIDs(userIDs)
		if err == nil {
			researchMetricsMap = metrics
		}
	}

	items := make([]*entity.AuthorEntity, len(authors))
	for i, author := range authors {
		items[i] = mapper.AuthorModelToEntity(&author.AuthorModel)
		items[i].ArticleCount = author.ArticleCount
		items[i].Photo = author.Photo

		if author.UserID != nil {
			userID := *author.UserID
			if user, ok := userMap[userID]; ok {
				if user.AcademicTitle != nil && *user.AcademicTitle != "" {
					title := *user.AcademicTitle
					items[i].AcademicTitle = &title
				}

				if string(user.AcademicDegree) != "" {
					if degree, exists := enum.GetDegreeByCode(string(user.AcademicDegree)); exists {
						label := degree.Label
						items[i].AcademicDegree = &label
					}
				}

				if user.ORCIDID != nil {
					items[i].ORCIDID = user.ORCIDID
				}
			}

			if job, ok := jobMap[userID]; ok {
				items[i].Job = mapper.JobModelToEntity(job)
			}

			if metrics, ok := researchMetricsMap[userID]; ok {
				items[i].ResearchMetrics = metrics
			}
		}
	}

	return items, nil
}

func (this *AuthorRepositoryImpl) GetOwnCoAuthorId(scienceID string) (uint, error) {
	var author models.AuthorModel
	err := this.db().Where("science_id = ?", scienceID).First(&author).Error
	if err != nil {
		return 0, infraErrors.Wrap(err)
	}
	return author.ID, nil
}

func (this *AuthorRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (r *AuthorRepositoryImpl) GetByScienceID(scienceID string) (*entity.AuthorEntity, error) {
	type AuthorWithUser struct {
		models.AuthorModel
		UserID *uint
		Photo  *string
	}

	var authorWithUser AuthorWithUser
	err := r.db().
		Table("authors a").
		Select("a.*, u.id as user_id, u.photo").
		Joins("LEFT JOIN users u ON u.science_id = a.science_id AND u.deleted_at IS NULL").
		Where("a.science_id = ? AND a.deleted_at IS NULL", scienceID).
		First(&authorWithUser).Error

	if err != nil {
		return nil, infraErrors.Wrap(err)
	}

	author := mapper.AuthorModelToEntity(&authorWithUser.AuthorModel)
	author.Photo = authorWithUser.Photo

	if authorWithUser.UserID != nil {
		userID := *authorWithUser.UserID

		var user models.UserModel
		err = r.db().Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error
		if err == nil {
			if user.AcademicTitle != nil && *user.AcademicTitle != "" {
				title := *user.AcademicTitle
				author.AcademicTitle = &title
			}

			if string(user.AcademicDegree) != "" {
				if degree, exists := enum.GetDegreeByCode(string(user.AcademicDegree)); exists {
					label := degree.Label
					author.AcademicDegree = &label
				}
			}

			if user.ORCIDID != nil {
				author.ORCIDID = user.ORCIDID
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infraErrors.Wrap(err)
		}

		var job models.JobModel
		err = r.db().
			Where("user_id = ? AND deleted_at IS NULL", userID).
			Order("created_at DESC").
			First(&job).Error
		if err == nil {
			author.Job = mapper.JobModelToEntity(&job)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infraErrors.Wrap(err)
		}

		researchMetrics, err := r.researchMetricRepo.GetByUserID(userID)
		if err == nil {
			author.ResearchMetrics = researchMetrics
		}
	}

	return author, nil
}

func (this *AuthorRepositoryImpl) SaveAuthor(author *entity.AuthorEntity) (*entity.AuthorEntity, error) {

	model := mapper.AuthorEntityToModel(author)
	result := this.db().Create(&model)
	if result.Error != nil {
		return nil, infraErrors.Wrap(result.Error)
	}
	return mapper.AuthorModelToEntity(&model), nil

}

func (this *AuthorRepositoryImpl) ExistingIds(ids []uint) ([]uint, error) {
	var existingIds []uint

	if err := this.db().Model(&models.AuthorModel{}).Where("id IN ?", ids).Pluck("id", &existingIds).Error; err != nil {
		return nil, err
	}

	return existingIds, nil
}

func (r *AuthorRepositoryImpl) GetAuthorsWithArticleCount(page, pageSize int, name, scienceID string, journalID *uint) (*entity.PagingEntity[entity.AuthorEntity], error) {
	type AuthorWithCount struct {
		models.AuthorModel
		ArticleCount int64
		UserID       *uint
		Photo        *string
	}

	var rawAuthors []AuthorWithCount

	query := r.db().
		Table("authors a").
		Select("a.*, COUNT(DISTINCT ac.article_model_id) as article_count, u.id as user_id, u.photo").
		Joins("JOIN article_co_authors ac ON a.id = ac.author_model_id").
		Joins("JOIN articles ar ON ac.article_model_id = ar.id AND ar.is_published = TRUE").
		Joins("LEFT JOIN users u ON u.science_id = a.science_id AND u.deleted_at IS NULL").
		Where("a.deleted_at IS NULL")

	if journalID != nil {
		query = query.Where("ar.journal_id = ?", *journalID)
	}

	if name != "" {
		query = query.Where("a.full_name ILIKE ?", "%"+name+"%")
	}
	if scienceID != "" {
		query = query.Where("a.science_id ILIKE ?", "%"+scienceID+"%")
	}

	query = query.
		Group("a.id, u.id, u.photo").
		Order("article_count DESC")

	if err := query.Scan(&rawAuthors).Error; err != nil {
		return nil, infraErrors.Wrap(err)
	}

	if len(rawAuthors) == 0 {
		return entity.NewPagingEntity(page, pageSize, 0, []*entity.AuthorEntity{}), nil
	}

	userIDs := make([]uint, 0, len(rawAuthors))
	authorIDMap := make(map[uint]int)
	uniqueAuthors := make([]AuthorWithCount, 0)
	seen := make(map[uint]bool)

	for _, a := range rawAuthors {
		if !seen[a.ID] {
			seen[a.ID] = true
			uniqueAuthors = append(uniqueAuthors, a)
		}
		if a.UserID != nil {
			userIDs = append(userIDs, *a.UserID)
			authorIDMap[a.ID] = len(uniqueAuthors) - 1
		}
	}

	total := int64(len(uniqueAuthors))

	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(uniqueAuthors) {
		return entity.NewPagingEntity(page, pageSize, total, []*entity.AuthorEntity{}), nil
	}
	if end > len(uniqueAuthors) {
		end = len(uniqueAuthors)
	}
	authorsPage := uniqueAuthors[start:end]

	userMap := make(map[uint]*models.UserModel)
	if len(userIDs) > 0 {
		var users []models.UserModel
		if err := r.db().Where("id IN ? AND deleted_at IS NULL", userIDs).Find(&users).Error; err == nil {
			for i := range users {
				userMap[users[i].ID] = &users[i]
			}
		}
	}

	jobMap := make(map[uint]*models.JobModel)
	if len(userIDs) > 0 {
		var jobs []models.JobModel
		sub := r.db().Table("jobs").
			Select("DISTINCT ON (user_id) *").
			Where("user_id IN ? AND deleted_at IS NULL", userIDs).
			Order("user_id, created_at DESC")
		_ = r.db().Table("(?) as j", sub).Find(&jobs).Error

		for i := range jobs {
			jobMap[jobs[i].UserID] = &jobs[i]
		}
	}

	researchMetricsMap := make(map[uint][]*entity.ResearchMetricEntity)
	if len(userIDs) > 0 {
		metrics, err := r.researchMetricRepo.GetByUserIDs(userIDs)
		if err == nil {
			researchMetricsMap = metrics
		}
	}

	items := make([]*entity.AuthorEntity, len(authorsPage))
	for i, a := range authorsPage {
		items[i] = mapper.AuthorModelToEntity(&a.AuthorModel)
		items[i].ArticleCount = a.ArticleCount
		items[i].Photo = a.Photo
		if a.UserID != nil {
			userID := *a.UserID
			if user, ok := userMap[userID]; ok {
				if user.AcademicTitle != nil && *user.AcademicTitle != "" {
					title := *user.AcademicTitle
					items[i].AcademicTitle = &title
				}

				if string(user.AcademicDegree) != "" {
					if degree, exists := enum.GetDegreeByCode(string(user.AcademicDegree)); exists {
						label := degree.Label
						items[i].AcademicDegree = &label
					}
				}

				if user.ORCIDID != nil {
					items[i].ORCIDID = user.ORCIDID
				}
			}

			if job, ok := jobMap[userID]; ok {
				items[i].Job = mapper.JobModelToEntity(job)
			}

			if metrics, ok := researchMetricsMap[userID]; ok {
				items[i].ResearchMetrics = metrics
			}
		}
	}

	return entity.NewPagingEntity(page, pageSize, total, items), nil
}

func (this *AuthorRepositoryImpl) GetJobs(authorIDs []uint) ([]*entity.JobWithAuthorIDEntity, error) {
	if len(authorIDs) == 0 {
		return []*entity.JobWithAuthorIDEntity{}, nil
	}

	// get the authors with their user_id through JOIN with users by science_id
	type AuthorWithUserID struct {
		AuthorID uint
		UserID   *uint
	}

	var authorsWithUsers []AuthorWithUserID
	err := this.db().
		Table("authors a").
		Select("a.id as author_id, u.id as user_id").
		Joins("LEFT JOIN users u ON u.science_id = a.science_id AND u.deleted_at IS NULL").
		Where("a.id IN ? AND a.deleted_at IS NULL", authorIDs).
		Scan(&authorsWithUsers).Error

	if err != nil {
		return nil, infraErrors.Wrap(err)
	}

	// create a map author_id -> user_id and collect unique user_ids
	authorToUserMap := make(map[uint]*uint)
	userIDSet := make(map[uint]bool)
	userIDs := make([]uint, 0)

	for _, au := range authorsWithUsers {
		if au.UserID != nil {
			authorToUserMap[au.AuthorID] = au.UserID
			userID := *au.UserID
			// add the user_id only if it doesn't exist
			if !userIDSet[userID] {
				userIDSet[userID] = true
				userIDs = append(userIDs, userID)
			}
		}
	}

	if len(userIDs) == 0 {
		return []*entity.JobWithAuthorIDEntity{}, nil
	}

	// get the last job for each user
	var jobs []models.JobModel
	sub := this.db().Table("jobs").
		Select("DISTINCT ON (user_id) *").
		Where("user_id IN ? AND deleted_at IS NULL", userIDs).
		Order("user_id, created_at DESC")

	err = this.db().Table("(?) as j", sub).Find(&jobs).Error
	if err != nil {
		return nil, infraErrors.Wrap(err)
	}

	// create a map user_id -> job for quick search
	userToJobMap := make(map[uint]*models.JobModel)
	for i := range jobs {
		userToJobMap[jobs[i].UserID] = &jobs[i]
	}

	// collect the result, matching the authors with their jobs
	result := make([]*entity.JobWithAuthorIDEntity, 0, len(authorIDs))
	for _, authorID := range authorIDs {
		if userID, ok := authorToUserMap[authorID]; ok {
			if job, jobOk := userToJobMap[*userID]; jobOk {
				result = append(result, entity.NewJobWithAuthorIDEntity(
					mapper.JobModelToEntity(job),
					authorID,
				))
			}
		}
	}

	return result, nil
}
