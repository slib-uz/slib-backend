package response

type KafkaStudyFieldResponse struct {
	ID       uint                     `json:"id"`
	Name     map[string]string        `json:"name"`
	Code     *uint                    `json:"code"`
	ParentID *uint                    `json:"parent_id"`
	Parent   *KafkaStudyFieldResponse `json:"parent"`
}

func NewKafkaStudyFieldResponse(id uint, name map[string]string, code, parentID *uint, parent *KafkaStudyFieldResponse) *KafkaStudyFieldResponse {
	return &KafkaStudyFieldResponse{ID: id, Name: name, Code: code, ParentID: parentID, Parent: parent}
}
