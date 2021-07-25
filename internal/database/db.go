package database

import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type DSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
