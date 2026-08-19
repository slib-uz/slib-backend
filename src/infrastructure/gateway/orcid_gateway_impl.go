package gateway

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"slib.uz/src/core/application/response"
	//"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
	. "slib.uz/src/infrastructure/gateway/response"
)

type ORCIDGatewayImpl struct {
	client       *network.CHTTpClient
	ClientID     string
	ClientSecret string
	BaseURL      string
}

// @inject
func NewORCIDGatewayImpl(client *network.CHTTpClient, env *config.Config) gateway.ORCIDGateway {
	return &ORCIDGatewayImpl{client: client, ClientID: env.ORCIDClientID, ClientSecret: env.ORCIDClientSecret, BaseURL: env.ORCIDBaseURL}
}

func (this *ORCIDGatewayImpl) AuthorizeUri(redirectUri string) string {
	return this.BaseURL + "/oauth/authorize?client_id=" + this.ClientID + "&response_type=code&scope=/authenticate&redirect_uri=" + redirectUri
}

func (this *ORCIDGatewayImpl) GetByCode(code string, redirectUri string) (string, error) {
	url := this.BaseURL + "/oauth/token"

	payload := map[string]string{
		"client_id":     this.ClientID,
		"client_secret": this.ClientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectUri,
	}

	resp, err := this.doRequest(url, payload)
	if err != nil {
		return "", fmt.Errorf("[ORCIDGatewayImpl] GetByCode error while making POST request: %w", err)
	}

	if resp.StatusCode != 200 {
		body, err := network.GetBody[map[string]any](resp)
		if err != nil {
			return "", fmt.Errorf("[ORCIDGatewayImpl] GetByCode received non-200 status code: %d, and error while parsing response body: %w", resp.StatusCode, err)
		}
		if resp.StatusCode == 400 {
			return "", response.NewOptionalResponse(200, response.CodeInvalidOrReusedCode, body, "")
		}
		return "", fmt.Errorf("[ORCIDGatewayImpl] GetByCode received non-200 status code: %d, response body: %v", resp.StatusCode, body)
	}

	body, err := network.GetBody[ORCIDTokenResponse](resp)
	if err != nil {
		return "", fmt.Errorf("[ORCIDGatewayImpl] GetByCode error while parsing response body: %w", err)
	}

	return body.Orcid, nil

}

func (this *ORCIDGatewayImpl) doRequest(urlStr string, form map[string]string) (*http.Response, error) {
	values := url.Values{}
	for k, v := range form {
		values.Set(k, v)
	}

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
