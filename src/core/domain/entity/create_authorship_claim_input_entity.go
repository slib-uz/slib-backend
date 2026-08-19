package entity

type CreateAuthorshipClaimInputEntity struct {
	UserID     uint
	ArticleIDs []uint
	Comment    string
}

func NewCreateAuthorshipClaimInputEntity(userID uint, articleIDs []uint, comment string) *CreateAuthorshipClaimInputEntity {
	return &CreateAuthorshipClaimInputEntity{
		UserID:     userID,
		ArticleIDs: articleIDs,
		Comment:    comment,
	}
}
