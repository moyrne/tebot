package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrRowsAffectedZero = errors.New("rows affected is zero")
)

func NewTransaction(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) (err error) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	tx, err := DB.Beginx()
	if err != nil {
		return errors.WithStack(err)
	}
	err = fn(ctx, tx)
	if err == nil {
		return errors.WithStack(tx.Commit())
	}
	if err := tx.Rollback(); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(err)
}
