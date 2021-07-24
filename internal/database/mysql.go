package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ConnectMySQL() (err error) {
	var dsnObj DSN
	if err := viper.UnmarshalKey("DB", &dsnObj); err != nil {
		return errors.WithStack(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", dsnObj.User, dsnObj.Password, dsnObj.Host, dsnObj.Port, dsnObj.DBName)
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return errors.WithStack(err)
	}
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(30)
	return errors.WithStack(DB.Ping())
}
