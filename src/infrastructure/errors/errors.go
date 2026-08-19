package errors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
)

func Wrap(e error) error {

	if e == nil {
		return nil
	}

	switch {
	case errors.Is(e, gorm.ErrRecordNotFound):
		return response.NotFoundError
	case isDuplicateKeyError(e):
		return response.ConflictError

	}
	return e
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505"
	}
	return false
}
