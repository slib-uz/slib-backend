package repository

type LegacyAuthorRepository interface {
	GetIDsByFullName(fullName string) ([]uint, error)
}
