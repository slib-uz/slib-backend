package entity

type UzSciApplicationAnswerEntity struct {
	FormID uint   `json:"form_id"`
	Value  string `json:"value"`
}

func NewUzSciApplicationAnswerEntity(formID uint, value string) *UzSciApplicationAnswerEntity {
	return &UzSciApplicationAnswerEntity{FormID: formID, Value: value}
}
