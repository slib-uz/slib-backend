package gateway

import (
	"fmt"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type AiDetectGatewayImpl struct {
	client       *network.CHTTpClient
	BaseURL      string
	ClientID     string
	ClientSecret string
}

// @inject
func NewAiDetectGateway(client *network.CHTTpClient, env *config.Config) gateway.AiDetectGateway {
	return &AiDetectGatewayImpl{
		client:       client,
		BaseURL:      env.AntiPlagBaseURL,
		ClientID:     env.AntiPlagClientID,
		ClientSecret: env.AntiPlagClientSecret,
	}
}

func (this *AiDetectGatewayImpl) Check(file []byte, fileName string) (*entity.AiDetectResultEntity, error) {
	url := this.BaseURL + "/api/ai-detects/check"

	payload := map[string]any{
		"file_name": fileName,
		"name":      fileName,
	}

	resp, err := this.client.PostFormData(url, file, fileName, payload, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to submit to ai-detect check: %w", err)
	}

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("AI detect response error:", data, err)
		return nil, fmt.Errorf("failed to submit to ai-detect check, status: %d", resp.StatusCode)
	}

	data, err := network.GetBody[response.AiDetectResultResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ai-detect response: %w", err)
	}

	return mapper.AiDetectResultResponseToEntity(data), nil
}

func (this *AiDetectGatewayImpl) GetResult(externalID uint) (*entity.AiDetectResultEntity, error) {
	url := fmt.Sprintf("%s/api/ai-detects/%d/result", this.BaseURL, externalID)

	resp, err := this.client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("gateway: failed to get ai-detect result: %w", err)
	}

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("AI detect response error (GetResult method):", data, err)
		return nil, fmt.Errorf("gateway: failed to get ai-detect result, status: %d", resp.StatusCode)
	}

	data, err := network.GetBody[response.AiDetectResultResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("gateway: failed to parse ai-detect result: %w", err)
	}

	return mapper.AiDetectResultResponseToEntity(data), nil
}
