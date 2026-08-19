package gateway

import (
	"fmt"
	"net/http"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/mapper"
	"slib.uz/src/infrastructure/gateway/network"
	. "slib.uz/src/infrastructure/gateway/response"
)

type ScienceIDGatewayImpl struct {
	Client       *network.CHTTpClient
	ClientID     string
	ClientSecret string
	BaseURL      string
}

// @inject
func NewScienceIDGateway(client *network.CHTTpClient, env *config.Config) gateway.ScienceIDGateway {
	return &ScienceIDGatewayImpl{
		Client:       client,
		ClientID:     env.ScienceIDClientID,
		ClientSecret: env.ScienceIDClientSecret,
		BaseURL:      env.ScienceIDBaseURL,
	}
}

func (this *ScienceIDGatewayImpl) GetUserByScienceID(scienceID string) (*entity.UserEntity, error) {
	url := fmt.Sprintf("%s/api/users/by-science-id/%s/", this.BaseURL, scienceID)

	resp, err := this.Client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, response.NotFoundError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ScienceIDGatewayImpl.GetUserByScienceID] unexpected status code: %d", resp.StatusCode)
	}

	userResponse, err := network.GetBody[UserDataResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("[ScienceIDGatewayImpl.GetUserByScienceID] unexpected response body: %v", err)
	}

	return mapper.UserDataResToEntity(userResponse), nil
}

func (this *ScienceIDGatewayImpl) GetAuthorByScienceID(scienceID string) (*entity.AuthorEntity, error) {
	url := this.BaseURL + "/api/shared/find-scientist/" + scienceID

	resp, err := this.Client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, response.NotFoundError
	}
	if resp.StatusCode == 200 {
		authorResponse, err := network.GetBody[AuthorResponse](resp)

		if err != nil {
			return nil, err
		}

		return mapper.AuthorResToEntity(authorResponse), nil
	}
	return nil, response.ExternalServiceError
}

func (this *ScienceIDGatewayImpl) GetUserScientificDataByScienceID(scienceID string) (*entity.AcademicDegreeEntity, *entity.AcademicTitleEntity, error) {
	url := this.BaseURL + "/profile/science-id/card/" + scienceID

	resp, err := this.Client.Get(url, network.NewBasicAuth(this.ClientID, this.ClientSecret), nil)

	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil, response.NotFoundError
	}
	if resp.StatusCode == 200 {
		scientificDataResponse, err := network.GetBody[ScientificDataResponse](resp)
		if err != nil {
			return nil, nil, err
		}

		return mapper.AcademicDegreeResToEntity(scientificDataResponse.AcademicDegree), mapper.AcademicTitleResToEntity(scientificDataResponse.AcademicTitle), nil
	}

	return nil, nil, response.ExternalServiceError
}
