package client

import (
	"github.com/hibiken/asynq"
	"slib.uz/src/infrastructure/config"
)

// @inject
func NewtAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddress})
}
