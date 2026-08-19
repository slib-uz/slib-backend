package db

import (
	"gorm.io/gorm"
)

type Database struct {
	GormDB *gorm.DB
}

// @inject
func NewDatabase(gormDB *gorm.DB) *Database {
	return &Database{GormDB: gormDB}
}
