package scheduler

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type SchedulerOption struct {
	Location *time.Location
	Queue    string
}

type SchedulerServer interface {
	Start()
	NewScheduler(cronSpec string, task *entity.AsyncTask, opt *SchedulerOption) error
}
