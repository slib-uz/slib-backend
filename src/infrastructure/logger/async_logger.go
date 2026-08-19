package logger

import "go.uber.org/zap"

type AsyncLogger struct {
	*zap.Logger
}

// @inject
func NewAsyncLogger() *AsyncLogger {
	return &AsyncLogger{Logger: NewDailyFileZapLogger("logs/async/")}
}
