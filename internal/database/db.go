package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/models"
	"github.com/pkg/errors"
	"log"
)

var DB *sqlx.DB

type DSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

var ErrDBNotConnect = errors.New("database not connected")

type Log struct{}

func (Log) Write(data []byte) (int, error) {
	if DB == nil {
		log.Println("database not connected")
		return 0, errors.WithStack(ErrDBNotConnect)
	}
	if err := NewTransaction(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return (&models.Log{Detail: string(data)}).Insert(ctx, tx)
	}); err != nil {
		return 0, err
	}
	return len(data), nil
}
