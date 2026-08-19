package gateway

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/gateway/response"
)

type SlibBillingGateway interface {
	GetJournalAccountBalance(journalID uint) (int64, error)
	GetUserAccountBalance(userID uint) (int64, error)
	//UseJournalServiceBalance(journalID uint, serviceCode enum.ServiceCode) error
	GetJournalServiceBalance(journalID uint, serviceCode enum.ServiceCode) (*response.SlibBillingServiceBalanceData, error)
	SendEvent(event *entity.OutboxEventEntity) error
}
