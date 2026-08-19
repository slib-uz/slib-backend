package gateway

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	core "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type UzSciGatewayImpl struct {
	Client  *network.CHTTpClient
	BaseURL string
}

// @inject
func NewUzSciGateway(client *network.CHTTpClient, env *config.Config) gateway.UzsciGateway {
	return &UzSciGatewayImpl{Client: client, BaseURL: env.UzSciBaseURL}
}

func (this *UzSciGatewayImpl) GetAllForms(ratingPeriodID uint) ([]*entity.UzSciFormEntity, error) {
	endpoint := this.BaseURL + "/api/v1/forms/all" + "?rating_period_id=" + strconv.FormatUint(uint64(ratingPeriodID), 10)

	resp, err := this.Client.Get(endpoint, nil, map[string]string{
		"accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetAllForms request error: %w", err)
	}

	if err := this.checkStatus(resp, http.StatusOK, "GetAllForms"); err != nil {
		return nil, err
	}

	body, err := network.GetBody[response.UzSciFormsResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetAllForms parse error: %w", err)
	}
	if !body.Success {
		return nil, core.NewFailResponse(502, "uzsci forms API returned success=false")
	}

	forms := make([]*entity.UzSciFormEntity, len(body.Data))
	for i := range body.Data {
		forms[i] = body.Data[i].ToEntity()
	}
	return forms, nil
}

func (this *UzSciGatewayImpl) GetPublicRatingPeriods(isActive *bool) ([]*entity.UzSciRatingPeriodEntity, error) {
	endpoint := this.BaseURL + "/api/v1/rating-periods/public"
	if isActive != nil {
		query := url.Values{}
		query.Set("is_active", strconv.FormatBool(*isActive))
		endpoint += "?" + query.Encode()
	}

	resp, err := this.Client.Get(endpoint, nil, map[string]string{
		"accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetPublicRatingPeriods request error: %w", err)
	}

	if err := this.checkStatus(resp, http.StatusOK, "GetPublicRatingPeriods"); err != nil {
		return nil, err
	}

	body, err := network.GetBody[response.UzSciRatingPeriodsResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetPublicRatingPeriods parse error: %w", err)
	}
	if !body.Success {
		return nil, core.NewFailResponse(502, "uzsci rating periods API returned success=false")
	}

	periods := make([]*entity.UzSciRatingPeriodEntity, len(body.Data))
	for i := range body.Data {
		periods[i] = body.Data[i].ToEntity()
	}
	return periods, nil
}

func (this *UzSciGatewayImpl) GetJournalByISSN(issnPaper, issnOnline string) (*response.UzSciJournalData, error) {
	query := url.Values{}
	if issnPaper != "" && utils.IsValidISSN(issnPaper) {
		query.Set("issn_paper", issnPaper)
	}
	if issnOnline != "" && utils.IsValidISSN(issnOnline) {
		query.Set("issn_online", issnOnline)
	}
	if len(query) == 0 {
		return nil, core.NewFailResponse(400, "at least one valid ISSN is required")
	}

	endpoint := this.BaseURL + "/api/v1/journals/by-issn?" + query.Encode()

	resp, err := this.Client.Get(endpoint, nil, map[string]string{
		"accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetJournalByISSN request error: %w", err)
	}

	if err := this.checkStatus(resp, http.StatusOK, "GetJournalByISSN"); err != nil {
		return nil, err
	}

	body, err := network.GetBody[response.UzSciJournalByISSNResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("[UzSciGateway] GetJournalByISSN parse error: %w", err)
	}
	if !body.Success || body.Data == nil {
		return nil, core.NewFailResponse(404, "uzsci journal not found")
	}

	return body.Data, nil
}

func (this *UzSciGatewayImpl) CreateApplication(periodID uint, journalID uint, answers []entity.UzSciApplicationAnswerEntity) error {
	endpoint := this.BaseURL + "/api/v1/rating-periods/" + strconv.FormatUint(uint64(periodID), 10) + "/applications"

	payloadAnswers := make([]response.UzSciCreateApplicationAnswerItem, len(answers))
	for i, answer := range answers {
		payloadAnswers[i] = response.UzSciCreateApplicationAnswerItem{
			FormID: answer.FormID,
			Value:  answer.Value,
		}
	}

	resp, err := this.Client.Post(endpoint, response.UzSciCreateApplicationPayload{
		JournalID: journalID,
		Answers:   payloadAnswers,
	}, nil, map[string]string{
		"accept": "application/json",
	})
	if err != nil {
		return fmt.Errorf("[UzSciGateway] CreateApplication request error: %w", err)
	}

	if err := this.checkStatus(resp, http.StatusCreated, "CreateApplication"); err != nil {
		return err
	}

	body, err := network.GetBody[response.UzSciCreateApplicationResponse](resp)
	if err != nil {
		return fmt.Errorf("[UzSciGateway] CreateApplication parse error: %w", err)
	}
	if !body.Success {
		return core.NewFailResponse(502, "uzsci create application API returned success=false")
	}

	return nil
}

func (this *UzSciGatewayImpl) checkStatus(resp *http.Response, expected int, op string) error {
	if resp.StatusCode == expected {
		return nil
	}

	body := network.GetErrorBody(resp)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return core.NotFoundError
	case http.StatusBadRequest:
		return core.NewFailResponse(400, fmt.Sprintf("uzsci %s: %v", op, body))
	case http.StatusConflict:
		return core.NewFailResponse(409, fmt.Sprintf("uzsci %s: %v", op, body))
	default:
		return fmt.Errorf("[UzSciGateway] %s unexpected status code: %d, body: %v", op, resp.StatusCode, body)
	}
}
