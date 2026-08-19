package entity

type EthicsPolicyEntity struct {
	ID      uint              `json:"id"`
	Content map[string]string `json:"content"`
}

func NewEthicsPolicyEntity(id uint, content map[string]string) *EthicsPolicyEntity {
	return &EthicsPolicyEntity{
		ID:      id,
		Content: content,
	}
}
