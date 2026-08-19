package entity

type ArticleFileEntity struct {
	PresignedURL string
	BinaryFile   []byte
}

// NewArticleFileEntity
func NewArticleFileEntity(presignedURL string, binaryFile []byte) *ArticleFileEntity {
	return &ArticleFileEntity{
		PresignedURL: presignedURL,
		BinaryFile:   binaryFile,
	}
}
