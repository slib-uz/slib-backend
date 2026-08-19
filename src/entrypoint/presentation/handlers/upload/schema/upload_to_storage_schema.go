package schema

type UploadToStorageRequest struct {
	FilePath string `json:"file_path" validate:"required"`
}
