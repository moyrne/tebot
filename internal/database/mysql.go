package database

import (
	"context"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var ErrRowsAffectedZero = errors.New("rows affected is zero")

func ConnectMySQL() (err error) {
	var dsnObj dbx.DSN
	if err := viper.UnmarshalKey("DB", &dsnObj); err != nil {
		return errors.WithStack(err)
	}

	DB, err = dbx.ConnectMySQL(dbx.DefaultMySQLDSN, dsnObj)
	return errors.WithStack(err)
}

func NewTransaction(ctx context.Context, fn func(ctx context.Context, tx dbx.Transaction) error) error {
	return dbx.NewTransaction(ctx, DB, fn)
}
