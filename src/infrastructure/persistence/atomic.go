package persistence

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/ports/session"
	db2 "slib.uz/src/infrastructure/db"
)

type AtomicTxImpl struct {
	db *db2.Database
}

// @inject
func NewAtomicTxImpl(db *db2.Database) session.Atomic {
	return &AtomicTxImpl{db: db}
}

func (this *AtomicTxImpl) Transaction(fn func(session.Tx) error) error {
	return this.db.GormDB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
