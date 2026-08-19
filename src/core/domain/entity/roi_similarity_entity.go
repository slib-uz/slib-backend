package entity

type RoiSimilarityEntity struct {
	Name    *map[string]string `json:"name"`
	ROI     *string            `json:"roi"`
	Percent float64            `json:"percent"`
	IsOwn   bool               `json:"is_own"`
}

func NewRoiSimilarityEntity(name *map[string]string, ROI *string, percent float64, isOwn bool) *RoiSimilarityEntity {
	return &RoiSimilarityEntity{Name: name, ROI: ROI, Percent: percent, IsOwn: isOwn}
}
