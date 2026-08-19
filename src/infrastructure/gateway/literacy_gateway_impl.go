package gateway

import (
	"encoding/base64"
	"fmt"
	_ "io/ioutil"
	"net/url"

	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type LiteracyGatewayImpl struct {
	Client  *network.CHTTpClient
	ApiKey  string
	BaseURL string
}

// @inject
func NewLiteracyGateway(client *network.CHTTpClient, env *config.Config) gateway.LiteracyGateway {
	return &LiteracyGatewayImpl{Client: client, ApiKey: env.LiteracyAPIKey, BaseURL: env.LiteracyBaseURL}
}

func (this *LiteracyGatewayImpl) SpellCheck(file []byte) ([]byte, error) {
	_url := this.BaseURL + "/api/document.php"

	postData := url.Values{
		"api_key":       {this.ApiKey},
		"task":          {"spellcheck"},
		"file_contents": {base64.StdEncoding.EncodeToString(file)},
	}

	resp, err := this.Client.PostForm(_url, postData)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	body, err := network.GetBody[response.SpellCheckResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	if !body.Success {
		return nil, fmt.Errorf("API error: %v", body)
	}

	outputData, _ := base64.StdEncoding.DecodeString(body.FileContents)

	return outputData, nil

}
