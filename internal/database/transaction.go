package database

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

var ErrRowsAffectedZero = errors.New("rows affected is zero")

func NewTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	// TODO 考虑是否 recover
	tx, err := DB.DB.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	err = fn(ctx, tx)
	if err == nil {
		return nil
	}
	if err := tx.Rollback(); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(err)
}
