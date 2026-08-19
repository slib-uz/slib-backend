package scheduler

import (
	"fmt"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	schedulerport "slib.uz/src/core/domain/ports/scheduler"
	"slib.uz/src/infrastructure/config"
)

type App struct {
	server schedulerport.SchedulerServer
	config *config.Config
}

// @inject
func NewApp(server schedulerport.SchedulerServer, config *config.Config) *App {
	return &App{server: server, config: config}
}

func (a *App) Init() {
	cronSpec := fmt.Sprintf("@every %dm", a.config.ViewsCountLifetimeMinute)

	a.must(a.server.NewScheduler(cronSpec, entity.NewAsyncTask(enum.TaskUpdateArticleViewsCount, nil), nil))
	a.must(a.server.NewScheduler(cronSpec, entity.NewAsyncTask(enum.TaskUpdateNewsViewsCount, nil), nil))
	a.must(a.server.NewScheduler(cronSpec, entity.NewAsyncTask(enum.TaskUpdateJournalViewsCount, nil), nil))

	utc5 := time.FixedZone("UTC+5", 5*60*60)
	a.must(a.server.NewScheduler("0 9 * * *", entity.NewAsyncTask(enum.TaskDeadlineReminder, nil), &schedulerport.SchedulerOption{
		Location: utc5,
	}))

	fmt.Println("[Scheduler] Scheduled tasks initialized successfully")
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) must(err error) {
	if err != nil {
		panic(err)
	}
}
