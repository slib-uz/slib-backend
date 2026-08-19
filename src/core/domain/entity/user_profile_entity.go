package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type UserProfileEntity struct {
	ID                      uint                     `json:"id"`
	UserID                  uint                     `json:"user_id"`
	FullName                string                   `json:"full_name"`
	ScienceID               string                   `json:"science_id"`
	Email                   string                   `json:"email"`
	Bio                     *string                  `json:"bio"`
	Photo                   *string                  `json:"photo"`
	Phone                   string                   `json:"phone"`
	Socials                 []*UserSocialEntity      `json:"socials"`
	Job                     *JobEntity               `json:"job"`
	ResearchMetrics         []*ResearchMetricEntity  `json:"research_metrics"`
	City                    *enum.CityCode           `json:"city" swaggertype:"string"`
	AcademicDegree          *enum.AcademicDegreeCode `json:"academic_degree" swaggertype:"string"`
	AcademicTitle           *string                  `json:"academic_title"`
	ORCIDID                 *string                  `json:"orcid_id"`
	PublishedArticlesCount  int64                    `json:"published_artilcles_count"`
	ProcessingArticlesCount int64                    `json:"proccesing_articles_count"`
}

func NewUserProfileEntity(id, userID uint, fullName string, scienceID string, email string, bio *string, photo *string, phone string, socials []*UserSocialEntity, job *JobEntity, researchMetrics []*ResearchMetricEntity, city *enum.CityCode, academicDegree *enum.AcademicDegreeCode, academicTitle *string, orcidID *string, publishedArticlesCount, processingArticlesCount int64) *UserProfileEntity {
	return &UserProfileEntity{
		ID:                      id,
		UserID:                  userID,
		FullName:                fullName,
		ScienceID:               scienceID,
		Email:                   email,
		Bio:                     bio,
		Photo:                   photo,
		Phone:                   phone,
		Socials:                 socials,
		Job:                     job,
		ResearchMetrics:         researchMetrics,
		City:                    city,
		AcademicDegree:          academicDegree,
		AcademicTitle:           academicTitle,
		ORCIDID:                 orcidID,
		PublishedArticlesCount:  publishedArticlesCount,
		ProcessingArticlesCount: processingArticlesCount,
	}
}
