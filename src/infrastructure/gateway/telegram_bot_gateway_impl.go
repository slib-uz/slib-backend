package gateway

import (
	"fmt"
	"net/url"

	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
)

type TelegramBotGatewayImpl struct {
	httpClient  *network.CHTTpClient
	botToken    string
	adminChatID string
}

// @inject
func NewTelegramBotGateway(env *config.Config, httpClient *network.CHTTpClient) gateway.TelegramBotGateway {
	return &TelegramBotGatewayImpl{
		httpClient:  httpClient,
		botToken:    env.TelegramBotToken,
		adminChatID: env.TelegramAdminChatID,
	}
}

func (this *TelegramBotGatewayImpl) AlertAdmin(message string) error {
	base := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown",
		url.PathEscape(this.botToken),
		url.PathEscape(this.adminChatID),
		url.QueryEscape(fmt.Sprintf("```\n%s\n```", message)),
	)

	resp, err := this.httpClient.Get(base, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send message to Telegram(TelegramBotGatewayImpl): %w", err)
	}
	if resp.StatusCode != 200 {
		body, err := network.GetBody[map[string]any](resp)
		if err != nil {
			return fmt.Errorf("failed to parse response body from Telegram(TelegramBotGatewayImpl): %w", err)
		}
		return fmt.Errorf("failed to send message to Telegram(TelegramBotGatewayImpl): %v status=%d", body, resp.StatusCode)
	}

	return nil
}
