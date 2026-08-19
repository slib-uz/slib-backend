package entity

type StudyFieldEntity struct {
	ID       uint                `json:"id"`
	Name     map[string]string   `json:"name"`
	ParentID *uint               `json:"parent_id"`
	Code     *uint               `json:"code"`
	Parent   *StudyFieldEntity   `json:"parent"`
	Children []*StudyFieldEntity `json:"children"`
}

func NewStudyFieldEntity(ID uint, name map[string]string, parentID *uint, code *uint, parent *StudyFieldEntity) *StudyFieldEntity {
	return &StudyFieldEntity{ID: ID, Name: name, ParentID: parentID, Code: code, Parent: parent}
}
