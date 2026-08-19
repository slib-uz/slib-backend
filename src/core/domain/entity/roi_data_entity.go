package entity

type ROIDataEntity struct {
	ID         uint           `json:"id"`
	RoiCode    string         `json:"roi_code"`
	ExternalID uint           `json:"external_id"`
	Data       *ArticleEntity `json:"data"`
	File       []byte         `json:"file"`
	Percent    float64        `json:"percent"`
}

func NewROIDataEntity(ID uint, roiCode string, externalID uint, data *ArticleEntity, file []byte, percent float64) *ROIDataEntity {
	return &ROIDataEntity{ID: ID, RoiCode: roiCode, ExternalID: externalID, Data: data, File: file, Percent: percent}
}
