package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ConnectPG() (err error) {
	var dsnObj DSN
	if err := viper.UnmarshalKey("DB", &dsnObj); err != nil {
		return errors.WithStack(err)
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dsnObj.Host,
		dsnObj.Port,
		dsnObj.User,
		dsnObj.Password,
		dsnObj.DBName)
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return errors.WithStack(err)
	}
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(30)
	return nil
}
