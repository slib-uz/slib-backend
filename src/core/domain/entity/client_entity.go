package entity

type ClientEntity struct {
	ID           uint
	ClientID     string
	ClientSecret string
	JournalID    *uint
	Name         string
	Description  string
	CallbackUrl  string
	IsActive     bool
}

func NewClientEntity(id uint, clientID, clientSecret, name, description, callbackUrl string, journalID *uint, isActive bool) *ClientEntity {
	return &ClientEntity{
		ID:           id,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		JournalID:    journalID,
		Name:         name,
		Description:  description,
		CallbackUrl:  callbackUrl,
		IsActive:     isActive,
	}
}

func NewClientBasicAuth(clientID, clientSecret string) *ClientEntity {
	return &ClientEntity{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}
