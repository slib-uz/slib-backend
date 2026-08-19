package response

type CoAuthorsResponse struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

type StudyFieldsResponse struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
	Code *uint             `json:"code"`
}

type ROIDetailResponse struct {
	ID              uint                   `json:"id"`
	Name            map[string]string      `json:"name"`
	PublicationDate string                 `json:"publication_date"`
	CoAuthors       []*CoAuthorsResponse   `json:"co_authors"`
	StudyFields     []*StudyFieldsResponse `json:"study_fields"`
	DOI             *string                `json:"doi,omitempty"`
	ROI             *string                `json:"roi,omitempty"`
	Percent         float64                `json:"percent"`
}
type RoiSuccessResponse struct {
	Status  int                `json:"status"`
	Code    *int               `json:"code"`
	Message string             `json:"message"`
	Data    *ROIDetailResponse `json:"data"`
}
