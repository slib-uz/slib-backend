package entity

type AuthorEntity struct {
	ID              uint                    `json:"id"`
	FullName        string                  `json:"full_name"`
	ScienceID       string                  `json:"science_id"`
	ArticleCount    int64                   `json:"article_count,omitempty"`
	Photo           *string                 `json:"photo,omitempty"`
	AcademicTitle   *string                 `json:"academic_title,omitempty"`
	AcademicDegree  *string                 `json:"academic_degree,omitempty"`
	Job             *JobEntity              `json:"job,omitempty"`
	User            *UserEntity             `json:"user,omitempty"`
	ORCIDID         *string                 `json:"orcid_id,omitempty"`
	ResearchMetrics []*ResearchMetricEntity `json:"research_metrics,omitempty"`
}

func NewAuthorEntity(ID uint, fullName string, scienceID string, articleCount int64, photo *string, academicTitle *string, academicDegree *string, job *JobEntity, orcidID *string, user *UserEntity) *AuthorEntity {
	return &AuthorEntity{
		ID:             ID,
		FullName:       fullName,
		ScienceID:      scienceID,
		ArticleCount:   articleCount,
		Photo:          photo,
		AcademicTitle:  academicTitle,
		AcademicDegree: academicDegree,
		Job:            job,
		User:           user,
		ORCIDID:        orcidID,
	}
}
