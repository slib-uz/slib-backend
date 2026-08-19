package scheduler

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

type SchedulerGroup struct {
	mu         sync.Mutex
	schedulers []*asynq.Scheduler
}

// @inject
func NewSchedulerGroup() *SchedulerGroup {
	return &SchedulerGroup{}
}

func (g *SchedulerGroup) Register(s *asynq.Scheduler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.schedulers = append(g.schedulers, s)
}

// Run barcha ro‘yxatdan o‘tgan scheduler‐larni parallel ishga tushiradi.
//   - SIGINT/SIGTERM kelganda ularni toza yopadi.
//   - Birontasi xato bersa, Run xatoni qaytaradi va qolganlarini ham yopadi.
func (g *SchedulerGroup) Run(ctx context.Context) error {
	if len(g.schedulers) == 0 {
		return nil // hech narsa yo‘q
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Har bir scheduler uchun gorutinda .Run() chaqiramiz
	for _, sch := range g.schedulers {
		s := sch
		eg.Go(func() error {
			return s.Run() // Run blokli; xato bo‘lsa errgroup o‘qiydi
		})
	}

	// OS signallarini tutib, graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			// errgroup konteksti bekor bo‘ldi (xato yoki boshqa gorutin tugadi)
		case <-sig:
			// tashqi signal keldi — shutdown
		}
		// barcha schedulerlarni to‘xtatamiz
		for _, s := range g.schedulers {
			s.Shutdown()
		}
		return nil
	})

	return eg.Wait()
}
