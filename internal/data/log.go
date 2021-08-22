package data

import (
	"context"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"time"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.LogRepo = logRepo{}

type logRepo struct{}

func NewLogRepo() cqhttp.LogRepo {
	return logRepo{}
}

func (l logRepo) Save(ctx context.Context, tx dbx.Transaction, log *cqhttp.Log) error {
	query := `insert into log (create_at,detail) values (?,?)`
	result, err := tx.ExecContext(ctx, query, time.Now(), log.Detail)
	if err != nil {
		return errors.WithStack(err)
	}
	log.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}
