package response

type PublishArticleResponse struct {
	ID      uint               `json:"id"`
	ROI     *string            `json:"roi"`
	Name    *map[string]string `json:"name"`
	Percent float64            `json:"percent"`
}

type UploadFileResponse struct {
	File string `json:"file"`
}

type ROIResponse struct {
	Status  int                     `json:"status"`
	Code    *int                    `json:"code"`
	Message string                  `json:"message"`
	Data    *PublishArticleResponse `json:"data"`
}
