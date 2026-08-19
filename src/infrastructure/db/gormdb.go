package db

import (
	"fmt"
	"time"

	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @inject
func NewGormDB(env *config.Config) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.DBSSLMode,
			),
		),
		&gorm.Config{
			Logger: logger.NewGormZapLogger(logger.GormLogConfig{
				SlowThreshold:        250 * time.Millisecond,
				IgnoreRecordNotFound: true,
			}),
		},
	)
	if err != nil {
		println("failed to connectUri persistence")
	}
	return db
}
