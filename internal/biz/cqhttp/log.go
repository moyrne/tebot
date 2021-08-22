package cqhttp

import (
	"context"
	"github.com/moyrne/tractor/dbx"
	"time"
)

type Log struct {
	ID       int64     `json:"id"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	Detail   string    `json:"detail"`
}

type LogRepo interface {
	Save(ctx context.Context, tx dbx.Transaction, log *Log) error
}
