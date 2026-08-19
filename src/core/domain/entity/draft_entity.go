package entity

type DraftEntity struct {
	ID     uint           `json:"id"`
	UserID uint           `json:"user_id"`
	Key    string         `json:"key"`
	Data   map[string]any `json:"data"`
}

func NewDraftEntity(ID uint, userID uint, key string, data map[string]any) *DraftEntity {
	return &DraftEntity{ID: ID, UserID: userID, Key: key, Data: data}
}
