package gateway

import (
	"fmt"
	"io"
	"path/filepath"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type ROIGatewayImpl struct {
	Client       *network.CHTTpClient
	ClientID     string
	ClientSecret string
	BaseURL      string
	storage      storage.FileStorage
}

// @inject
func NewROIGateway(
	client *network.CHTTpClient,
	env *config.Config,
	storage storage.FileStorage,
) gateway.ROIGateway {
	return &ROIGatewayImpl{
		Client:       client,
		ClientID:     env.ROIClientID,
		ClientSecret: env.ROIClientSecret,
		BaseURL:      env.ROIBaseURL,
		storage:      storage,
	}
}

func (this *ROIGatewayImpl) Publish(data *entity.ArticleEntity, sourceLink string) (*entity.ROIDataEntity, error) {
	fileBytes, fErr := this.storage.GetObject(enum.BucketPrivate, data.ContentFile)
	if fErr != nil {
		return nil, fmt.Errorf("[ROIGateway] Get File Object: %w", fErr)
	}

	fileName := filepath.Base(data.ContentFile)
	uploadedPath, uErr := this.uploadFile(fileBytes, fileName)
	if uErr != nil {
		return nil, uErr
	}

	payload := mapper.ArticleEntityToPayload(data, uploadedPath, sourceLink)
	url := this.BaseURL + "/api/article/publish"

	resp, err := this.Client.Post(url, payload, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("[ROIGateway] Publish Post request: %w", err)
	}

	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		body, parseErr := network.GetBody[response.ROIResponse](resp)
		if parseErr != nil {
			return nil, fmt.Errorf("[ROIGateway] Parse Body: %w", parseErr)
		}
		return entity.NewROIDataEntity(data.ID, *body.Data.ROI, body.Data.ID, data, fileBytes, body.Data.Percent), nil
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("[ROIGateway] Status Code: Article not found")
	}

	return nil, fmt.Errorf("[ROIGateway] Publish: unexpected status code %d", resp.StatusCode)
}

func (this *ROIGatewayImpl) uploadFile(fileBytes []byte, filename string) (*string, error) {
	url := this.BaseURL + "/api/storage/upload"

	// Send file as form-data
	resp, err := this.Client.PostFormData(
		url, fileBytes,
		filename, map[string]any{},
		network.NewBasicAuth(this.ClientID, this.ClientSecret),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[ROIGateway] UploadFile: %w", err)
	}

	// Handle successful response
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		result, parseErr := network.GetBody[response.UploadFileResponse](resp)
		if parseErr != nil {
			return nil, fmt.Errorf("[ROIGateway] UploadFile: %w", parseErr)
		}

		return &result.File, nil
	}

	// Handle error responses
	_, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("[ROIGateway] UploadFile read error response: %w", readErr)
	}

	return nil, fmt.Errorf("[ROIGateway] UploadFile: unexpected status code %d", resp.StatusCode)
}

func (this *ROIGatewayImpl) UpdateArticle(article *entity.ROIArticleUpdateEntity) error {
	if article.File != "" {
		fileBytes, fErr := this.storage.GetObject(enum.BucketPrivate, article.File)
		if fErr != nil {
			return fmt.Errorf("[ROIGateway] UpdateArticle: %w", fErr)
		}

		fileName := filepath.Base(article.File)
		uploadedPath, uErr := this.uploadFile(fileBytes, fileName)
		if uErr != nil {
			return uErr
		}

		article.File = *uploadedPath
	}

	url := fmt.Sprintf("%s/api/article/update", this.BaseURL)

	resp, err := this.Client.Put(url, article, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return fmt.Errorf("[ROIGateway] UpdateArticle error while making PUT request: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("[ROIGateway] UpdateArticle unexpected status code %d", resp.StatusCode)
	}

	return nil
}
