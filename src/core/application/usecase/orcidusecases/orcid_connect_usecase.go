package orcidusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type OrcidConnectUseCase struct {
	gateway    gateway.ORCIDGateway
	repository repository.UserRepository
}

// @inject
func NewOrcidConnectUseCase(gateway gateway.ORCIDGateway, repository repository.UserRepository) *OrcidConnectUseCase {
	return &OrcidConnectUseCase{gateway: gateway, repository: repository}
}

func (this *OrcidConnectUseCase) Execute(requestUserID uint, code string, redirectUri string) error {
	orcid, err := this.gateway.GetByCode(code, redirectUri)
	if err != nil {
		return err
	}

	user, err := this.repository.GetByOrcid(orcid)
	if err != nil && !errors.Is(err, response.NotFoundError) {
		return err
	}

	if user != nil {
		if user.ID == requestUserID {
			return nil
		}
		return response.NewOptionalResponse(200, response.CodeOrcidAlreadyConnected, map[string]any{
			"user_id":    user.ID,
			"orcid":      user.ORCIDID,
			"full_name":  user.FullName,
			"science_id": user.ScienceID,
		}, "")
	}

	return this.repository.UpdateOrcid(requestUserID, orcid)
}
