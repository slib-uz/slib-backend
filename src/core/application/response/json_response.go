package response

type JsonResponse struct {
	Status   int     `json:"status"`
	Code     *int    `json:"code"`
	CodeName *string `json:"_code"`
	Message  *string `json:"message"`
	Data     any     `json:"data"`
}
