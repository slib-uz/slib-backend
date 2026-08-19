package cache

type ArticleViewsCountCache interface {
	Add(userKey string, articleID uint) (int64, error)
	GetAll() (map[uint]int64, error)
}
