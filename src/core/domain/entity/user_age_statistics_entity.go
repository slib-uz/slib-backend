package entity

// UserAgeStatisticsEntity represents age statistics for users
type UserAgeStatisticsEntity struct {
	Total     int64 `json:"total"`
	Age0_17   int64 `json:"0-17"`
	Age18_24  int64 `json:"18-24"`
	Age25_34  int64 `json:"25-34"`
	Age35_44  int64 `json:"35-44"`
	Age45_54  int64 `json:"45-54"`
	Age55_64  int64 `json:"55-64"`
	Age65Plus int64 `json:"65+"`
}

func NewUserAgeStatisticsEntity(total, age0_17, age18_24, age25_34, age35_44, age45_54, age55_64, age65Plus int64) *UserAgeStatisticsEntity {
	return &UserAgeStatisticsEntity{
		Total:     total,
		Age0_17:   age0_17,
		Age18_24:  age18_24,
		Age25_34:  age25_34,
		Age35_44:  age35_44,
		Age45_54:  age45_54,
		Age55_64:  age55_64,
		Age65Plus: age65Plus,
	}
}
