package entity

type GuideRetrieveEntity struct {
	ID          uint              `json:"id"`
	Title       map[string]string `json:"title"`
	Description map[string]string `json:"description"`
	FilePath    string            `json:"file_path"`
	VideoUrl    string            `json:"video_url"`
}

func NewGuideRetrieveEntity(ID uint, title map[string]string, description map[string]string, filePath string, videoUrl string) *GuideRetrieveEntity {
	return &GuideRetrieveEntity{ID: ID, Title: title, Description: description, FilePath: filePath, VideoUrl: videoUrl}
}
