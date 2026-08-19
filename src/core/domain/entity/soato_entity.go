package entity

type SoatoClassifierEntity struct {
	ID       uint                     `json:"id"`
	Code     uint                     `json:"code"`
	Name     map[string]string        `json:"name"`
	ParentID *uint                    `json:"parent_id"`
	Parent   *SoatoClassifierEntity   `json:"parent,omitempty"`
	Children []*SoatoClassifierEntity `json:"children,omitempty"`
}
