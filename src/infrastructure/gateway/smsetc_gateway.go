package gateway

import (
	"fmt"

	"github.com/google/uuid"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
)

type SmsEtcGateway struct {
	client    *network.CHTTpClient
	BaseURL   string
	Username  string
	Password  string
	AlphaName string
}

// @inject
func NewSmsEtcGateway(client *network.CHTTpClient, env *config.Config) gateway.SmsGateway {
	return &SmsEtcGateway{client: client, BaseURL: env.SmsEtcBaseURL, Username: env.SmsEtcUsername, Password: env.SmsEtcPassword, AlphaName: env.SmsAlphaName}
}

func (this *SmsEtcGateway) Send(phone string, message string) error {
	url := fmt.Sprintf("%s/single-sms", this.BaseURL)

	payload := map[string]any{
		"header": map[string]string{"login": this.Username, "pwd": this.Password, "CgPN": this.AlphaName},
		"body":   map[string]string{"message_id_in": uuid.NewString(), "CdPN": phone, "text": message},
	}

	resp, err := this.client.Post(url, payload, nil, nil)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("[SmsEtcGateway] Send error while making POST request: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("[SmsEtcGateway] Send unexpected status code %d, response body: %v", resp.StatusCode, network.GetErrorBody(resp))
	}
	return nil
}
