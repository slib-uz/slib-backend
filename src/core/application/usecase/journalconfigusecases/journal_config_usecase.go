package journalconfigusecases

import (
	"errors"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalConfigUseCase struct {
	repository     repository.JournalConfigRepository
	billingGateway gateway.SlibBillingGateway
}

// @inject
func NewJournalConfigUseCase(
	repository repository.JournalConfigRepository,
	billingGateway gateway.SlibBillingGateway,
) *JournalConfigUseCase {
	return &JournalConfigUseCase{
		repository:     repository,
		billingGateway: billingGateway,
	}
}

func (uc *JournalConfigUseCase) Execute(journalID uint, websiteUrl string, user *entity.UserBasicEntity) (*entity.JournalConfigEntity, error) {
	if journalID != 0 {
		return uc.getByJournalID(journalID, user)
	}

	return uc.getByWebsiteURL(websiteUrl, user)
}

func (uc *JournalConfigUseCase) getByJournalID(journalID uint, user *entity.UserBasicEntity) (*entity.JournalConfigEntity, error) {
	if uc.allowCheckBalance(user) {
		if err := uc.checkBalance(journalID); err != nil {
			return nil, err
		}
	}

	journalConfig, err := uc.repository.GetByJournalID(journalID)
	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return nil, response.NewOptionalResponse(200, response.CodeNotFound, nil, "journal config not found")
		}
		return nil, err
	}

	return journalConfig, nil
}

func (uc *JournalConfigUseCase) getByWebsiteURL(websiteUrl string, user *entity.UserBasicEntity) (*entity.JournalConfigEntity, error) {
	journalConfig, err := uc.repository.GetByWebsiteURL(websiteUrl)
	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return nil, response.NewOptionalResponse(200, response.CodeNotFound, nil, "journal config not found")
		}
		return nil, err
	}

	if uc.allowCheckBalance(user) {
		if err := uc.checkBalance(journalConfig.JournalID); err != nil {
			return nil, err
		}

		if err := uc.checkStatus(journalConfig); err != nil {
			return nil, err
		}
	}

	return journalConfig, nil
}

func (uc *JournalConfigUseCase) checkBalance(journalID uint) error {
	balance, err := uc.billingGateway.GetJournalServiceBalance(journalID, enum.ServiceJournalWebsite)
	if err != nil {
		if errors.Is(err, response.NotFound) {
			return err
		}

		return response.NewOptionalResponse(200, response.CodeIntegrationError, nil, "integration error")
	}

	if balance.ExpiresAt == nil || balance.ExpiresAt.Before(time.Now()) {
		return response.NewOptionalResponse(200, response.CodeInsufficientBalance, nil, "insufficient balance")
	}

	return nil
}

func (this *JournalConfigUseCase) checkStatus(journalConfig *entity.JournalConfigEntity) error {
	if journalConfig == nil || !journalConfig.IsActive {
		return response.NewOptionalResponse(200, response.CodeNotFound, nil, "journal config not found")
	}

	return nil
}

func (uc *JournalConfigUseCase) List(creatorID, journalID uint, isActive *bool, page, pageSize int) (*entity.PagingEntity[entity.JournalConfigEntity], error) {
	return uc.repository.List(creatorID, journalID, isActive, page, pageSize)
}

// adminka uchun check balance ishlamasligi kerak, websayt uchun ishlashi kerak.
func (uc *JournalConfigUseCase) allowCheckBalance(user *entity.UserBasicEntity) bool {
	if user == nil {
		return true
	}

	return false
}
