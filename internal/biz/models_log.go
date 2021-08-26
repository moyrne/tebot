package biz

import (
	"context"
	"time"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Log struct {
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	Detail   string    `json:"detail"`
}

type LogRepo interface {
	Save(ctx context.Context, tx dbx.Transaction, log *Log) error
	Write([]byte) (int, error)
}

type DBHook struct {
	db   dbx.Database
	repo LogRepo
}

func NewDBHook(db dbx.Database, repo LogRepo) *DBHook {
	return &DBHook{db: db, repo: repo}
}

func (d *DBHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (d *DBHook) Fire(e *logrus.Entry) error {
	serialized, err := e.Logger.Formatter.Format(e)
	if err != nil {
		return errors.WithStack(err)
	}
	return dbx.NewTransaction(context.Background(), d.db, func(ctx context.Context, tx dbx.Transaction) error {
		return d.repo.Save(context.Background(), tx, &Log{
			CreateAt: time.Now(),
			Detail:   string(serialized),
		})
	})
}
