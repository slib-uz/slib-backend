package response

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

// POST /api/v1/journals/{journal_id}/use-service-balance

type SlibBillingUseBalanceResponse struct {
	Status  int                       `json:"status"`
	Success bool                      `json:"success"`
	Code    int                       `json:"code,omitempty"`
	Message string                    `json:"message,omitempty"`
	Data    SlibBillingUseBalanceData `json:"data,omitempty"`
}

type SlibBillingUseBalanceData struct {
	Message string `json:"message"`
}

// GET /api/v1/journals/{journal_id}/service-balance?service_code={service_code}

type SlibBillingServiceBalanceResponse struct {
	Status  int                            `json:"status"`
	Success bool                           `json:"success"`
	Code    int                            `json:"code,omitempty"`
	Message string                         `json:"message,omitempty"`
	Data    *SlibBillingServiceBalanceData `json:"data,omitempty"`
}

type SlibBillingServiceBalanceData struct {
	AccountID int                              `json:"account_id"`
	ServiceID int                              `json:"service_id"`
	Service   SlibBillingServiceBalanceService `json:"service"`
	Available int                              `json:"available"`
	Used      int                              `json:"used"`
	ExpiresAt *time.Time                       `json:"expires_at,omitempty"`
}

func (this *SlibBillingServiceBalanceData) IsAvailable() bool {
	if this.Service.Type == enum.ServiceTimeBased {
		expiresAt := this.ExpiresAt
		if expiresAt == nil {
			return false
		}

		return time.Until(*this.ExpiresAt) > 0
	}

	return this.Available > 0
}

type SlibBillingServiceBalanceService struct {
	ID          int                             `json:"id"`
	Code        string                          `json:"code"`
	Name        SlibBillingServiceBalanceLocale `json:"name"`
	Description SlibBillingServiceBalanceLocale `json:"description"`
	Type        enum.ServiceType                `json:"type"`
	IsActive    bool                            `json:"is_active"`
}

type SlibBillingServiceBalanceLocale struct {
	En string `json:"en"`
	Ru string `json:"ru"`
	Uz string `json:"uz"`
}
