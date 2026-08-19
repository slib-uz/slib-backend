package entity

type PublicOfferEntity struct {
	ID          uint              `json:"id"`
	Description map[string]string `json:"description"`
}

func NewPublicOfferEntity(ID uint, description map[string]string) *PublicOfferEntity {
	return &PublicOfferEntity{ID: ID, Description: description}
}
