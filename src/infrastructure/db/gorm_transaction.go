package db

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/ports/unitofwork"
)

type GormTransaction struct {
	db *Database
}

// @inject
func NewGormTransaction(db *Database) unitofwork.Atomic {
	return &GormTransaction{db: db}
}

func (this *GormTransaction) Transaction(fn func(tx unitofwork.Tx) error) error {
	return this.db.GormDB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
