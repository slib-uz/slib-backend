package logger

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogConfig struct {
	SlowThreshold        time.Duration
	IgnoreRecordNotFound bool
}

type GormZapLogger struct {
	cfg GormLogConfig
}

func NewGormZapLogger(cfg GormLogConfig) *GormZapLogger {
	if cfg.SlowThreshold == 0 {
		cfg.SlowThreshold = 250 * time.Millisecond
	}

	return &GormZapLogger{cfg: cfg}
}

func (this *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return this
}

func (this *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{})  {}
func (this *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{})  {}
func (this *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {}

func (this *GormZapLogger) loggerFrom(ctx context.Context) *zap.Logger {
	if ctx != nil {
		if lg := FromContext(ctx); lg != nil {
			return lg
		}
	}
	return nil
}

func (this *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	lg := this.loggerFrom(ctx)
	if lg == nil {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && this.cfg.IgnoreRecordNotFound && errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	fields := []zap.Field{
		zap.Duration("duration", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
		zap.String("scope", "sql"),
	}

	if err != nil {
		lg.Error("sql error", append(fields, zap.Error(err))...)
		return
	}

	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	isSelect := strings.HasPrefix(sqlUpper, "SELECT")
	isSlow := this.cfg.SlowThreshold > 0 && elapsed >= this.cfg.SlowThreshold

	if isSlow || !isSelect {
		if isSlow {
			lg.Warn("sql slow", fields...)
		} else {
			lg.Info("sql write", fields...)
		}
	}
}
