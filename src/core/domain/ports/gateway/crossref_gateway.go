package gateway

import "slib.uz/src/core/domain/entity"

type CrossRefDepositParams struct {
	BatchID     string
	Username    string
	Password    string
	JournalName string
	Journal     *entity.JournalEntity
	Articles    []*entity.ArticleEntity
}

type CrossRefGateway interface {
	Deposit(params *CrossRefDepositParams) (*entity.CrossRefDepositResultEntity, error)
	CheckAuth(username, password string) (bool, error)
	GetPrefixInfo(prefix string) (*entity.CrossRefPrefixInfoEntity, error)
}
