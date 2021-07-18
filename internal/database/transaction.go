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
	var tx *sqlx.Tx
	tx, err = DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "beginx")
	}
	// recover
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = errors.WithMessagef(err, "rollback %v", e)
			}
			return
		}
		err = errors.Wrap(tx.Commit(), "tx commit")
	}()

	return fn(ctx, tx)
}
