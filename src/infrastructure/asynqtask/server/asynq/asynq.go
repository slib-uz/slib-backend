package asynq

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/config"
)

// @inject
func NewAsynqServer(redisClient *asynq.RedisClientOpt, config asynq.Config) *asynq.Server {
	return asynq.NewServer(redisClient, config)
}

// @inject
func NewAsynqServerMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

// @inject
func NewRedisClientOpt0(cfg *config.Config) *asynq.RedisClientOpt {
	return &asynq.RedisClientOpt{Addr: cfg.RedisAddress, DB: 0}
}

// @inject
func NewAsynqMon(redisClient *asynq.RedisClientOpt) *asynqmon.HTTPHandler {
	return asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring/tasks",
		RedisConnOpt: redisClient,
	})
}

// @inject
func NewConfig() asynq.Config {
	return asynq.Config{

		Concurrency: 10,
		Queues: map[string]int{
			"critical": 5,
			"default":  4,
			"low":      1,
		},
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			switch t.Type() {
			case string(enum.TaskAntiPlagStatusUpdate):
				return 2 * time.Minute
			default:
				return min(1800*time.Second, 3*time.Second<<uint(n-1))
			}
		},
	}
}
