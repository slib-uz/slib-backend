package response

type ScientificDataResponse struct {
	AcademicDegree *AcademicDegreeResponse `json:"academic_degree"`
	AcademicTitle  *AcademicTitleResponse  `json:"academic_title"`
}

func NewScientificDataResponse(academicDegree *AcademicDegreeResponse, academicTitle *AcademicTitleResponse) *ScientificDataResponse {
	return &ScientificDataResponse{
		AcademicDegree: academicDegree,
		AcademicTitle:  academicTitle,
	}
}
