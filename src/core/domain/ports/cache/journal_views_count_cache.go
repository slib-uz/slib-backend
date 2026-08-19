package cache

type JournalViewsCountCache interface {
	Add(userKey string, journalID uint) (int64, error)
	GetAll() (map[uint]int64, error)
}
