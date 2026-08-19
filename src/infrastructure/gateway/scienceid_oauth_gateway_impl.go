package gateway

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	coreErr "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type ScienceIdOAuthGatewayImpl struct {
	Client       *network.CHTTpClient
	ClientID     string
	ClientSecret string
	BaseURL      string
	CodeVerifier string
}

// @inject
func NewScienceIdOAuthGateway(client *network.CHTTpClient, cfg *config.Config) gateway.ScienceIdOAuthGateway {
	return &ScienceIdOAuthGatewayImpl{
		Client:       client,
		ClientID:     cfg.ScienceIDClientID,
		ClientSecret: cfg.ScienceIDClientSecret,
		BaseURL:      cfg.ScienceIDBaseURL,
		CodeVerifier: cfg.OAuthCodeVerifier,
	}
}

func (this *ScienceIdOAuthGatewayImpl) GetUserByCode(code string) (*entity.UserEntity, error) {

	url := fmt.Sprintf("%s/oauth/user/get/", this.BaseURL)

	payload := map[string]any{
		"code":          code,
		"code_verifier": this.CodeVerifier,
	}

	resp, err := this.Client.Post(url, payload, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)

	if err != nil {
		return nil, err
	}

	defer network.Closer(resp.Body)()

	if resp.StatusCode == 400 {
		var errorResponse map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}

		if code, ok := errorResponse["code"].(float64); ok && int(code) == 4201 {
			return nil, coreErr.InvalidCodeError
		}
	}

	if resp.StatusCode == 200 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(body))

		var user response.UserResponse

		err = json.Unmarshal(body, &user)
		if err != nil {
			return nil, err

		}

		return mapper.UserDataResToEntity(&user.Data), nil
	}
	fmt.Println("ErrorData response:", resp.StatusCode, err)
	return nil, coreErr.ExternalServiceError

}

func (this *ScienceIdOAuthGatewayImpl) AuthorizeURL(redirectURL string) string {
	codeChallenge := this.generateCodeChallenge(this.CodeVerifier)

	return fmt.Sprintf(
		"%s/oauth/authorize?client_id=%s&redirect_uri=%s&code_challenge=%s&state=%d",
		this.BaseURL, this.ClientID, redirectURL, codeChallenge, rand.Int(),
	)
}

func (this *ScienceIdOAuthGatewayImpl) generateCodeChallenge(codeVerifier string) string {
	sha256Hash := sha256.Sum256([]byte(codeVerifier))

	codeChallenge := base64.URLEncoding.EncodeToString(sha256Hash[:])

	return codeChallenge[:len(codeChallenge)-1]
}
