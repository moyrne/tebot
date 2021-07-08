package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func ConnectPG() (err error) {
	dsn, err := loadDSN()
	if err != nil {
		return err
	}
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "sqlx.Connect")
	}
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(30)
	return nil
}

type PgDSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func loadDSN() (string, error) {
	var dsn PgDSN
	if err := viper.UnmarshalKey("DB", &dsn); err != nil {
		return "", errors.Wrap(err, "unmarshal key DB")
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dsn.Host,
		dsn.Port,
		dsn.User,
		dsn.Password,
		dsn.DBName), nil
}
