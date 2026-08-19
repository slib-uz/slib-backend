package entity

// UserGenderStatisticsEntity represents gender statistics for users
type UserGenderStatisticsEntity struct {
	Total  int64 `json:"total"`
	Male   int64 `json:"male"`
	Female int64 `json:"female"`
}

func NewUserGenderStatisticsEntity(total, male, female int64) *UserGenderStatisticsEntity {
	return &UserGenderStatisticsEntity{
		Total:  total,
		Male:   male,
		Female: female,
	}
}
