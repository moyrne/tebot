package data

import (
	"context"
	"time"

	"github.com/moyrne/tebot/internal/biz"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ biz.LogRepo = logRepo{}

type logRepo struct{}

func NewLogRepo() biz.LogRepo {
	return logRepo{}
}

func (l logRepo) Save(ctx context.Context, tx dbx.Transaction, log *biz.Log) error {
	query := `insert into log (create_at,detail) values (?,?)`
	result, err := tx.ExecContext(ctx, query, time.Now(), log.Detail)
	if err != nil {
		return errors.WithStack(err)
	}
	log.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (l logRepo) Write(data []byte) (int, error) {
	err := database.NewTransaction(context.Background(), func(ctx context.Context, tx dbx.Transaction) error {
		return l.Save(ctx, tx, &biz.Log{
			CreateAt: time.Now(),
			Detail:   string(data),
		})
	})
	return len(data), err
}
