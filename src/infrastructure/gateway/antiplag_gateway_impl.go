package gateway

import (
	"fmt"
	"io"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type AntiPlagGatewayImpl struct {
	client       *network.CHTTpClient
	BaseURL      string
	ClientID     string
	ClientSecret string
}

// @inject
func NewAntiPlagGateway(client *network.CHTTpClient, env *config.Config) gateway.AntiPlagGateway {
	return &AntiPlagGatewayImpl{client: client, BaseURL: env.AntiPlagBaseURL, ClientID: env.AntiPlagClientID, ClientSecret: env.AntiPlagClientSecret}
}

func (this *AntiPlagGatewayImpl) Check(file []byte, fileName string) (*entity.AntiPlagResultEntity, error) {
	url := this.BaseURL + "/api/documents/check"

	payload := map[string]any{
		"file_name":   fileName,
		"name":        fileName,
		"category_id": 2,
	}

	resp, err := this.client.PostFormData(url, file, fileName, payload, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to submit to antiplag check: %w", err)
	}

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("Antiplag response error:", data, err)
		return nil, fmt.Errorf("failed to submit to antiplag check, status: %d", resp.StatusCode)
	}

	data, err := network.GetBody[response.AntiPlagResultResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse antiplag response: %w", err)
	}

	return mapper.AntiPlagResultResponseToEntity(data), nil

}

func (this *AntiPlagGatewayImpl) GetStatus(externalID uint) int {
	//TODO implement me
	panic("implement me")
}

func (this *AntiPlagGatewayImpl) GetResult(externalID uint) (*entity.AntiPlagResultEntity, error) {

	url := fmt.Sprintf("%s/api/documents/%d/result", this.BaseURL, externalID)

	resp, err := this.client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("gateway: failed to submit to antiplag result: %w", err)
	}

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("Antiplag response error (GetResult method):", data, err)
		return nil, fmt.Errorf("gateway: failed to submit to antiplag result, status: %d", resp.StatusCode)
	}

	data, err := network.GetBody[response.AntiPlagResultResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("gateway: failed to parse antiplag result: %w", err)
	}

	return mapper.AntiPlagResultResponseToEntity(data), nil
}

func (this *AntiPlagGatewayImpl) GetBalance() (int, error) {
	url := this.BaseURL + "/api/balance"

	resp, err := this.client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return 0, fmt.Errorf("gateway: failed to get antiplag balance: %w", err)
	}

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("Antiplag response error (GetBalance method):", data, err)
		return 0, fmt.Errorf("gateway: failed to get antiplag balance, status: %d", resp.StatusCode)
	}

	data, err := network.GetBody[response.AntiPlagBalanceResponse](resp)
	if err != nil {
		return 0, fmt.Errorf("gateway: failed to parse antiplag balance: %w", err)
	}

	return data.Balance, nil
}

func (this *AntiPlagGatewayImpl) GenerateCertificate(externalID uint, authorFullName, fileName string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/documents/%d/certificate", this.BaseURL, externalID)

	body := map[string]any{
		"author_full_name": authorFullName,
		"name":             fileName,
	}

	resp, err := this.client.Post(url, body, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("gateway: failed to generate antiplag certificate(AntiPlagGatewayImpl): %w", err)
	}

	defer utils.Closer(resp.Body)()

	if resp.StatusCode != 200 {
		data, err := network.GetBody[map[string]any](resp)
		fmt.Println("Antiplag response error (GenerateCertificate method AntiPlagGatewayImpl):", data, err)
		return nil, fmt.Errorf("(AntiPlagGatewayImpl) gateway: failed to generate antiplag certificate, status: %d", resp.StatusCode)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("(AntiPlagGatewayImpl) gateway: failed to read antiplag certificate response body: %w", err)
	}

	return fileBytes, nil
}
