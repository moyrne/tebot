package models

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Log struct {
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	Detail   string    `json:"detail"`
}

func (l *Log) Insert(ctx context.Context, tx *sqlx.Tx) error {
	query := `insert into log (create_at,detail) values (?,?)`
	result, err := tx.ExecContext(ctx, query, time.Now(), l.Detail)
	if err != nil {
		return errors.WithStack(err)
	}
	l.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}
