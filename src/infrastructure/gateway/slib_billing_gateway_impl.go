package gateway

import (
	"fmt"
	"net/http"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
	. "slib.uz/src/infrastructure/gateway/response"
)

type SlibBillingGatewayImpl struct {
	Client   *network.CHTTpClient
	BaseURL  string
	Username string
	Secret   string
}

// @inject
func NewSlibBillingGateway(client *network.CHTTpClient, cfg *config.Config) gateway.SlibBillingGateway {
	return &SlibBillingGatewayImpl{
		Client:   client,
		BaseURL:  cfg.SlibBillingBaseURL,
		Username: cfg.ClientBasicAuthUsername,
		Secret:   cfg.ClientBasicAuthSecret,
	}
}

func (this *SlibBillingGatewayImpl) GetJournalAccountBalance(journalID uint) (int64, error) {
	url := fmt.Sprintf("%s/api/v1/journals/%d/account", this.BaseURL, journalID)

	resp, err := this.Client.Get(url, network.NewBasicAuth(this.Username, this.Secret), nil)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return 0, response.NotFoundError
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalAccountBalance] unexpected status code: %d", resp.StatusCode)
	}

	billingResponse, err := network.GetBody[SlibBillingAccountResponse](resp)
	if err != nil {
		return 0, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalAccountBalance] unexpected response body: %v", err)
	}

	return billingResponse.Data.Balance, nil
}

func (this *SlibBillingGatewayImpl) GetUserAccountBalance(userID uint) (int64, error) {
	// url := fmt.Sprintf("%s/api/v1/users/%d/account", this.BaseURL, userID)

	// resp, err := this.Client.Get(url, network.NewBasicAuth(this.Username, this.Secret), nil)
	// if err != nil {
	// 	return 0, err
	// }

	// if resp.StatusCode == http.StatusNotFound {
	// 	return 0, response.NotFoundError
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	return 0, fmt.Errorf("[SlibBillingGatewayImpl.GetUserAccountBalance] unexpected status code: %d", resp.StatusCode)
	// }

	// billingResponse, err := network.GetBody[SlibBillingAccountResponse](resp)
	// if err != nil {
	// 	return 0, fmt.Errorf("[SlibBillingGatewayImpl.GetUserAccountBalance] unexpected response body: %v", err)
	// }

	// --------------------------------------

	// 	IMPLEMENT ME. FIX SLIB BILLING API

	// --------------------------------------

	return 0, nil
}

//func (this *SlibBillingGatewayImpl) UseJournalServiceBalance(journalID uint, serviceCode enum.ServiceCode) error {
//	url := fmt.Sprintf("%s/api/v1/journals/%d/use-service-balance?service_code=%s", this.BaseURL, journalID, serviceCode.String())
//
//	resp, err := this.Client.Post(url, nil, network.NewBasicAuth(this.Username, this.Secret), nil)
//	if err != nil {
//		return fmt.Errorf("[SlibBillingGatewayImpl.UseJournalServiceBalance] request failed: %w", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("[SlibBillingGatewayImpl.UseJournalServiceBalance] unexpected status code: %d", resp.StatusCode)
//	}
//
//	billingResponse, err := network.GetBody[SlibBillingUseBalanceResponse](resp)
//	if err != nil {
//		return fmt.Errorf("[SlibBillingGatewayImpl.UseJournalServiceBalance] unexpected response body: %v", err)
//	}
//
//	if !billingResponse.Success {
//		if billingResponse.Code == response.CodeInsufficientBalance.Code {
//			return response.InsufficientBalance
//		}
//		return fmt.Errorf("[SlibBillingGatewayImpl.UseJournalServiceBalance] billing service error: %s", billingResponse.Message)
//	}
//
//	return nil
//}

func (this *SlibBillingGatewayImpl) SendEvent(event *entity.OutboxEventEntity) error {
	url := fmt.Sprintf("%s/api/v1/inbox-events", this.BaseURL)

	body := map[string]any{
		"event_id": event.EventID,
		"version":  event.Version,
		"event":    event.Event,
		"payload":  event.Payload,
	}

	resp, err := this.Client.Post(url, body, network.NewBasicAuth(this.Username, this.Secret), nil)
	if err != nil {
		return fmt.Errorf("[SlibBillingGatewayImpl.SendEvent] request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[SlibBillingGatewayImpl.SendEvent] unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (this *SlibBillingGatewayImpl) GetJournalServiceBalance(journalID uint, serviceCode enum.ServiceCode) (*SlibBillingServiceBalanceData, error) {
	url := fmt.Sprintf("%s/api/v1/journals/service-balance?service_code=%s", this.BaseURL, serviceCode.String())

	resp, err := this.Client.Get(url, network.NewBasicAuth(this.Username, this.Secret), map[string]string{
		"Journal": fmt.Sprintf("%d", journalID),
	})
	if err != nil {
		return nil, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalServiceBalance] request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalServiceBalance] unexpected status code: %d", resp.StatusCode)
	}

	billingResponse, err := network.GetBody[SlibBillingServiceBalanceResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalServiceBalance] unexpected response body: %v", err)
	}

	if !billingResponse.Success {
		if billingResponse.Code == response.CodeNotFound.Code {
			return nil, response.NotFound
		}

		return nil, fmt.Errorf("[SlibBillingGatewayImpl.GetJournalServiceBalance] billing service error: %s", billingResponse.Message)
	}

	return billingResponse.Data, nil
}
