package repository

import (
	"errors"
	"fmt"
	"sort"

	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/session"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type UserRepositoryImpl struct {
	database          *db.Database
	profileRepository repository.UserProfileRepository
	roleRepository    repository.UserRoleRepository
}

// @inject
func NewUserRepository(database *db.Database, profileRepository repository.UserProfileRepository, roleRepository repository.UserRoleRepository) repository.UserRepository {
	return &UserRepositoryImpl{
		database:          database,
		profileRepository: profileRepository,
		roleRepository:    roleRepository,
	}

}

func (this *UserRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *UserRepositoryImpl) GetByPhoneNumber(phoneNumber string) (*entity.UserEntity, error) {
	var _model models.UserModel
	if err := this.db().Where("phone_number = ?", phoneNumber).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.UserModelToEntity(&_model), nil
}

func (this *UserRepositoryImpl) CreateByPhoneNumber(tx session.Tx, user *entity.UserEntity) (uint, error) {
	_db := tx.(*gorm.DB)

	var _model = mapper.UserEntityToModel(user)

	if err := _db.Create(_model).Error; err != nil {
		return 0, _errors.Wrap(err)
	}
	return _model.ID, nil
}

func (this *UserRepositoryImpl) UpdateOrcid(userID uint, orcid string) error {
	err := this.db().Model(&models.UserModel{}).Where("id = ?", userID).Update("orc_id_id", orcid).Error
	return _errors.Wrap(err)
}

func (this *UserRepositoryImpl) GetByOrcid(orcid string) (*entity.UserEntity, error) {
	var _model models.UserModel
	if err := this.db().Where("orc_id_id = ?", orcid).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.UserModelToEntity(&_model), nil
}

func (this *UserRepositoryImpl) Save(user *entity.UserEntity) (*entity.UserEntity, error) {

	var userFromDB, err = this.GetByPin(*user.Pin)
	if err != nil {
		return nil, _errors.Wrap(err)
	}

	var modelFromEntity = mapper.UserEntityToModel(user)

	if userFromDB == nil {
		err = this.db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(modelFromEntity).Error; err != nil {
				return _errors.Wrap(err)
			}
			return this.profileRepository.Create(tx, modelFromEntity.ID)
		})

		return mapper.UserModelToEntity(modelFromEntity), nil
	}
	// TODO: update only fields that are allowed to update
	if err := this.db().Model(userFromDB).Updates(modelFromEntity).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.UserModelToEntity(userFromDB), nil
}

func (this *UserRepositoryImpl) GetById(id uint) (*entity.UserEntity, error) {

	var user models.UserModel

	result := this.db().
		Preload("Roles").
		Preload("Roles.Journal").
		Preload("Roles.Publisher").
		Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, _errors.Wrap(result.Error)
	}
	return mapper.UserModelToEntity(&user), nil
}

func (this *UserRepositoryImpl) GetByPin(pin string) (*models.UserModel, error) {
	var user models.UserModel
	result := this.db().Where("pin = ?", pin).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (this *UserRepositoryImpl) GetByScienceId(scienceId string) (*entity.UserEntity, error) {
	var user models.UserModel

	if err := this.db().Preload("Roles").Where("science_id = ?", scienceId).First(&user).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.UserModelToEntity(&user), nil

}

func (this *UserRepositoryImpl) FindByPhone(phone string) (*entity.UserEntity, error, bool) {
	var user models.UserModel
	if err := this.db().Where("phone_number = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _errors.Wrap(err), true
		} else {
			return nil, _errors.Wrap(err), false
		}
	}
	return mapper.UserModelToEntity(&user), nil, false
}

func (this *UserRepositoryImpl) Update(id uint, user *entity.UserEntity) (*entity.UserEntity, error) {
	// TODO: update only fields that are allowed to update
	var userModel models.UserModel

	if err := this.db().Where("id = ?", id).First(&userModel).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	updateMap := make(map[string]interface{})

	if user.Photo != "" && user.Photo != userModel.Photo {
		updateMap["photo"] = user.Photo
	}

	if user.Email != "" && user.Email != userModel.Email {
		updateMap["email"] = user.Email
	}

	if string(user.AcademicDegree) != "" && string(user.AcademicDegree) != string(userModel.AcademicDegree) {
		updateMap["academic_degree"] = user.AcademicDegree
	}

	if user.AcademicTitle != nil {
		if userModel.AcademicTitle == nil || *user.AcademicTitle != *userModel.AcademicTitle {
			updateMap["academic_title"] = user.AcademicTitle
		}
	}

	if user.ORCIDID != nil && user.ORCIDID != userModel.ORCIDID {
		updateMap["orc_id_id"] = user.ORCIDID
	}

	if len(updateMap) > 0 {
		if err := this.db().Model(&userModel).Updates(updateMap).Error; err != nil {
			return nil, _errors.Wrap(err)
		}
	}

	if err := this.db().Preload("Roles").Where("id = ?", id).First(&userModel).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.UserModelToEntity(&userModel), nil
}

func (this *UserRepositoryImpl) GetGenderStatistics() (*entity.UserGenderStatisticsEntity, error) {
	type genderStats struct {
		Total  int64
		Male   int64
		Female int64
	}

	var stats genderStats

	err := this.db().Table("users").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN gender = 1 THEN 1 END) as male,
			COUNT(CASE WHEN gender = 2 THEN 1 END) as female
		`).
		Where("deleted_at IS NULL").
		Scan(&stats).Error

	if err != nil {
		return nil, _errors.Wrap(err)
	}

	return entity.NewUserGenderStatisticsEntity(stats.Total, stats.Male, stats.Female), nil
}

func (this *UserRepositoryImpl) GetAgeStatistics() (*entity.UserAgeStatisticsEntity, error) {
	type ageStats struct {
		Total     int64 `json:"total"`
		Age0_17   int64 `json:"age0_17"`
		Age18_24  int64 `json:"age18_24"`
		Age25_34  int64 `json:"age25_34"`
		Age35_44  int64 `json:"age35_44"`
		Age45_54  int64 `json:"age45_54"`
		Age55_64  int64 `json:"age55_64"`
		Age65Plus int64 `json:"age65plus"`
	}

	var stats ageStats

	err := this.db().Table("users").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 0 AND 17 THEN 1 END) as age0_17,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 18 AND 24 THEN 1 END) as age18_24,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 25 AND 34 THEN 1 END) as age25_34,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 35 AND 44 THEN 1 END) as age35_44,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 45 AND 54 THEN 1 END) as age45_54,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) BETWEEN 55 AND 64 THEN 1 END) as age55_64,
			COUNT(CASE WHEN birth_date IS NOT NULL AND EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) >= 65 THEN 1 END) as age65plus
		`).
		Where("deleted_at IS NULL").
		Scan(&stats).Error

	if err != nil {
		return nil, _errors.Wrap(err)
	}

	return entity.NewUserAgeStatisticsEntity(
		stats.Total,
		stats.Age0_17,
		stats.Age18_24,
		stats.Age25_34,
		stats.Age35_44,
		stats.Age45_54,
		stats.Age55_64,
		stats.Age65Plus,
	), nil
}

func (this *UserRepositoryImpl) GetArticleStatisticsByYear(userID uint, year int) (*entity.UserArticleStatisticsByYearEntity, error) {
	type monthStats struct {
		Month        int   `json:"month"`
		Articles     int64 `json:"articles"`
		Applications int64 `json:"applications"`
	}

	var stats []monthStats

	err := this.db().Raw(`
		WITH months AS (
			SELECT generate_series(1, 12) AS month
		),
		user_articles AS (
			SELECT
				EXTRACT(MONTH FROM a.publication_date)::int AS month,
				COUNT(DISTINCT a.id) AS article_count
			FROM articles a
			JOIN article_author_affiliations aaa ON aaa.article_id = a.id
			JOIN authors au ON au.id = aaa.author_id
			JOIN users u ON u.science_id = au.science_id
			WHERE
				u.id = ?
				AND a.deleted_at IS NULL
				AND a.is_published = true
				AND EXTRACT(YEAR FROM a.publication_date) = ?
			GROUP BY EXTRACT(MONTH FROM a.publication_date)
		),
		user_applications AS (
			SELECT
				EXTRACT(MONTH FROM aa.created_at)::int AS month,
				COUNT(DISTINCT aa.id) AS application_count
			FROM article_applications aa
			WHERE
				aa.user_id = ?
				AND aa.deleted_at IS NULL
				AND EXTRACT(YEAR FROM aa.created_at) = ?
			GROUP BY EXTRACT(MONTH FROM aa.created_at)
		)
		SELECT
			m.month,
			COALESCE(ua.article_count, 0) AS articles,
			COALESCE(uapp.application_count, 0) AS applications
		FROM months m
		LEFT JOIN user_articles ua ON ua.month = m.month
		LEFT JOIN user_applications uapp ON uapp.month = m.month
		ORDER BY m.month
	`, userID, year, userID, year).Scan(&stats).Error

	if err != nil {
		return nil, _errors.Wrap(err)
	}

	months := make(map[string]entity.MonthStatisticsEntity)
	var totalArticles int64
	var totalApplications int64

	for _, stat := range stats {
		monthKey := fmt.Sprintf("%d", stat.Month)
		months[monthKey] = entity.NewMonthStatisticsEntity(stat.Articles, stat.Applications)
		totalArticles += stat.Articles
		totalApplications += stat.Applications
	}

	total := entity.NewMonthStatisticsEntity(totalArticles, totalApplications)

	return entity.NewUserArticleStatisticsByYearEntity(total, months), nil
}

func (this *UserRepositoryImpl) GetArticleStatisticsByYearRange(userID uint, fromYear, toYear int) (*entity.UserArticleStatisticsByYearRangeEntity, error) {
	type yearStats struct {
		Year         int   `json:"year"`
		Articles     int64 `json:"articles"`
		Applications int64 `json:"applications"`
	}

	var stats []yearStats

	err := this.db().Raw(`
		WITH years AS (
			SELECT generate_series(?::int, ?::int) AS year
		),
		user_articles AS (
			SELECT
				EXTRACT(YEAR FROM a.publication_date)::int AS year,
				COUNT(DISTINCT a.id) AS article_count
			FROM articles a
			JOIN article_author_affiliations aaa ON aaa.article_id = a.id
			JOIN authors au ON au.id = aaa.author_id
			JOIN users u ON u.science_id = au.science_id
			WHERE
				u.id = ?
				AND a.deleted_at IS NULL
				AND a.is_published = true
				AND EXTRACT(YEAR FROM a.publication_date) BETWEEN ? AND ?
			GROUP BY EXTRACT(YEAR FROM a.publication_date)
		),
		user_applications AS (
			SELECT
				EXTRACT(YEAR FROM aa.created_at)::int AS year,
				COUNT(DISTINCT aa.id) AS application_count
			FROM article_applications aa
			WHERE
				aa.user_id = ?
				AND aa.deleted_at IS NULL
				AND EXTRACT(YEAR FROM aa.created_at) BETWEEN ? AND ?
			GROUP BY EXTRACT(YEAR FROM aa.created_at)
		)
		SELECT
			y.year,
			COALESCE(ua.article_count, 0) AS articles,
			COALESCE(uapp.application_count, 0) AS applications
		FROM years y
		LEFT JOIN user_articles ua ON ua.year = y.year
		LEFT JOIN user_applications uapp ON uapp.year = y.year
		ORDER BY y.year
	`, fromYear, toYear, userID, fromYear, toYear, userID, fromYear, toYear).Scan(&stats).Error

	if err != nil {
		return nil, _errors.Wrap(err)
	}

	years := make([]entity.YearStatisticsItemEntity, 0, len(stats))
	var totalArticles int64
	var totalApplications int64

	for _, stat := range stats {
		yearKey := fmt.Sprintf("%d", stat.Year)
		years = append(years, entity.NewYearStatisticsItemEntity(yearKey, stat.Articles, stat.Applications))
		totalArticles += stat.Articles
		totalApplications += stat.Applications
	}

	sort.Slice(years, func(i, j int) bool {
		return years[i].Year < years[j].Year
	})

	total := entity.NewMonthStatisticsEntity(totalArticles, totalApplications)

	return entity.NewUserArticleStatisticsByYearRangeEntity(total, years), nil
}

func (this *UserRepositoryImpl) GetUserActivityStartYear(userID uint) (int, error) {
	type activityYear struct {
		FirstApplicationYear *int `json:"first_application_year"`
		FirstArticleYear     *int `json:"first_article_year"`
	}

	var result activityYear

	err := this.db().Raw(`
		SELECT
			(SELECT EXTRACT(YEAR FROM MIN(aa.created_at))::int
			 FROM article_applications aa
			 WHERE aa.user_id = ? AND aa.deleted_at IS NULL) AS first_application_year,
			(SELECT EXTRACT(YEAR FROM MIN(a.publication_date))::int
			 FROM articles a
			 JOIN article_author_affiliations aaa ON aaa.article_id = a.id
			 JOIN authors au ON au.id = aaa.author_id
			 JOIN users u ON u.science_id = au.science_id
			 WHERE u.id = ? AND a.deleted_at IS NULL AND a.is_published = true) AS first_article_year
	`, userID, userID).Scan(&result).Error

	if err != nil {
		return 0, _errors.Wrap(err)
	}

	var earliestYear *int

	if result.FirstApplicationYear != nil && result.FirstArticleYear != nil {
		if *result.FirstApplicationYear < *result.FirstArticleYear {
			earliestYear = result.FirstApplicationYear
		} else {
			earliestYear = result.FirstArticleYear
		}
	} else if result.FirstApplicationYear != nil {
		earliestYear = result.FirstApplicationYear
	} else if result.FirstArticleYear != nil {
		earliestYear = result.FirstArticleYear
	}

	if earliestYear == nil {
		return 0, nil
	}

	return *earliestYear, nil
}

func (this *UserRepositoryImpl) GetAll(page, pageSize int, search string) (*entity.PagingEntity[entity.UserEntity], error) {
	var users []*models.UserModel
	var total int64

	query := this.db().Model(&models.UserModel{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR pin ILIKE ? OR email ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := query.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("Roles").
		Preload("Roles.Journal").
		Preload("Roles.Publisher").
		Find(&users).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	// Fetch Article Counts efficiently
	var usersEntities []*entity.UserEntity
	if len(users) > 0 {
		var scienceIDs []string
		userMap := make(map[string]*entity.UserEntity)

		for _, user := range users {
			_entity := mapper.UserModelToEntity(user)
			usersEntities = append(usersEntities, _entity)
			if user.ScienceID != "" {
				scienceIDs = append(scienceIDs, user.ScienceID)
				userMap[user.ScienceID] = _entity
			}
		}

		if len(scienceIDs) > 0 {
			type Result struct {
				ScienceID string
				Count     int64
			}
			var results []Result

			this.db().
				Table("authors au").
				Select("au.science_id, COUNT(DISTINCT a.id) as count").
				Joins("JOIN article_co_authors ac ON au.id = ac.author_model_id").
				Joins("JOIN articles a ON ac.article_model_id = a.id").
				Where("au.science_id IN (?) AND a.deleted_at IS NULL AND a.is_published = true", scienceIDs).
				Group("au.science_id").
				Scan(&results)

			for _, res := range results {
				if user, ok := userMap[res.ScienceID]; ok {
					user.ArticleCount = res.Count
				}
			}
		}
	}

	return entity.NewPagingEntity(
		page,
		pageSize,
		total,
		usersEntities,
	), nil
}

func (this *UserRepositoryImpl) GetDetailByID(id uint) (*entity.UserDetailEntity, error) {
	var user models.UserModel

	if err := this.db().
		Preload("Roles").
		Preload("Roles.Journal").
		Preload("Roles.Publisher").
		Where("id = ?", id).First(&user).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	userMe := mapper.UserMeProfileModelToEntity(&models.UserProfileModel{
		UserID: user.ID,
		User:   &user,
	})

	var job *entity.JobEntity
	var jobModel models.JobModel
	if err := this.db().
		Where("user_id = ? AND deleted_at IS NULL", id).
		Order("created_at DESC").
		First(&jobModel).Error; err == nil {
		job = mapper.JobModelToEntity(&jobModel)
	}

	var researchMetrics []*entity.ResearchMetricEntity
	var matricModels []*models.ResearchMetricModel
	if err := this.db().Where("user_id = ?", id).Find(&matricModels).Error; err == nil {
		researchMetrics = mapper.ResearchMetricListModelToEntity(matricModels)
	}

	var articleCount int64
	if user.ScienceID != "" {
		this.db().
			Table("articles a").
			Joins("JOIN article_co_authors ac ON a.id = ac.article_model_id").
			Joins("JOIN authors au ON ac.author_model_id = au.id").
			Where("au.science_id = ? AND a.deleted_at IS NULL AND a.is_published = true", user.ScienceID).
			Count(&articleCount)
	}

	var academicDegree *string
	if string(user.AcademicDegree) != "" {
		if degree, exists := enum.GetDegreeByCode(string(user.AcademicDegree)); exists {
			label := degree.Label
			academicDegree = &label
		}
	}

	return entity.NewUserDetailEntity(
		userMe,
		job,
		researchMetrics,
		articleCount,
		academicDegree,
		user.AcademicTitle,
		user.ORCIDID,
	), nil
}
