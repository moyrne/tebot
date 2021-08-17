package data

import (
	"context"
	"time"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

type Log struct {
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	Detail   string    `json:"detail"`
}

func (l *Log) Insert(ctx context.Context, tx dbx.Transaction) error {
	query := `insert into log (create_at,detail) values (?,?)`
	result, err := tx.ExecContext(ctx, query, time.Now(), l.Detail)
	if err != nil {
		return errors.WithStack(err)
	}
	l.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}
