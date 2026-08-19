package logger

import "go.uber.org/zap"

type HttpLogger struct {
	*zap.Logger
}

// @inject
func NewHttpLogger() *HttpLogger {
	return &HttpLogger{Logger: NewDailyFileZapLogger("logs/http/")}
}
