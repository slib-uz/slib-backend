package response

type KafkaJournalResponse struct {
	ID         uint              `json:"id"`
	Name       map[string]string `json:"name"`
	ShortName  map[string]string `json:"short_name,omitempty"`
	ISSNPaper  *string           `json:"issn_paper,omitempty"`
	ISSNOnline *string           `json:"issn_online,omitempty"`
}

func NewKafkaJournalResponse(id uint, name map[string]string, shortName map[string]string, issnPaper, issnOnline *string) *KafkaJournalResponse {
	return &KafkaJournalResponse{ID: id, Name: name, ShortName: shortName, ISSNPaper: issnPaper, ISSNOnline: issnOnline}
}
