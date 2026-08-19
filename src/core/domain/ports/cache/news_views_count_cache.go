package cache

type NewsViewsCountCache interface {
	Add(userKey string, newsID uint) (int64, error)
	GetAll() (map[uint]int64, error)
}
