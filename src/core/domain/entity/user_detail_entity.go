package entity

type UserDetailEntity struct {
	UserMeEntity
	Job             *JobEntity              `json:"job"`
	ResearchMetrics []*ResearchMetricEntity `json:"research_metrics"`
	ArticleCount    int64                   `json:"article_count"`
	AcademicDegree  *string                 `json:"academic_degree"`
	AcademicTitle   *string                 `json:"academic_title"`
	ORCIDID         *string                 `json:"orcid_id"`
}

func NewUserDetailEntity(
	userMe *UserMeEntity,
	job *JobEntity,
	researchMetrics []*ResearchMetricEntity,
	articleCount int64,
	academicDegree *string,
	academicTitle *string,
	orcidID *string,
) *UserDetailEntity {
	return &UserDetailEntity{
		UserMeEntity:    *userMe,
		Job:             job,
		ResearchMetrics: researchMetrics,
		ArticleCount:    articleCount,
		AcademicDegree:  academicDegree,
		AcademicTitle:   academicTitle,
		ORCIDID:         orcidID,
	}
}
