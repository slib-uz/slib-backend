package scheduler

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/entity"
	schedulerport "slib.uz/src/core/domain/ports/scheduler"
	"slib.uz/src/infrastructure/config"
)

type SchedulerServerImpl struct {
	redisClientOpt *asynq.RedisClientOpt
	group          *SchedulerGroup
}

// @inject
func NewSchedulerServerImpl(env *config.Config, group *SchedulerGroup) schedulerport.SchedulerServer {
	redisClientOpt := &asynq.RedisClientOpt{Addr: env.RedisAddress, DB: 0}
	return &SchedulerServerImpl{redisClientOpt: redisClientOpt, group: group}
}

func (s *SchedulerServerImpl) Start() {
	if err := s.group.Run(context.Background()); err != nil {
		panic(err)
	}
}

func (s *SchedulerServerImpl) NewScheduler(cronSpec string, task *entity.AsyncTask, opt *schedulerport.SchedulerOption) error {
	var schedulerOpts *asynq.SchedulerOpts
	if opt != nil && opt.Location != nil {
		schedulerOpts = &asynq.SchedulerOpts{Location: opt.Location}
	}

	sch := asynq.NewScheduler(s.redisClientOpt, schedulerOpts)

	queue := "low"
	if opt != nil && opt.Queue != "" {
		queue = opt.Queue
	}

	asynqTask := asynq.NewTask(string(task.TaskType), task.Payload)
	if _, err := sch.Register(cronSpec, asynqTask, asynq.Queue(queue)); err != nil {
		return fmt.Errorf("[SchedulerServerImpl] error registering scheduler: %v", err)
	}

	s.group.Register(sch)
	return nil
}
