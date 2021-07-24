package database

import (
	_ "unsafe"

	"context"
	"database/sql"
	. "github.com/agiledragon/gomonkey"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// pg_sql.(*Tx).Rollback 被内联, 无法打桩, 需要对私有方法打桩
//go:linkname rollbackFn database/pg_sql.(*Tx).rollback
func rollbackFn(*sql.Tx, bool) error

func TestNewTransaction(t *testing.T) {
	var rollback, committed bool
	reset := func() {
		rollback = false
		committed = false
	}
	ApplyMethod(reflect.TypeOf(&sqlx.DB{}), "Beginx", func(db *sqlx.DB) (*sqlx.Tx, error) {
		return &sqlx.Tx{}, nil
	})
	txType := reflect.TypeOf(&sql.Tx{})
	ApplyFunc(rollbackFn, func(tx *sql.Tx, discardConn bool) error {
		rollback = true
		return nil
	})
	ApplyMethod(txType, "Commit", func(tx *sql.Tx) error {
		committed = true
		return nil
	})

	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		defer reset()
		assert.Equal(t, nil, NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			return nil
		}))
		assert.Equal(t, false, rollback)
		assert.Equal(t, true, committed)
	})
	t.Run("return error", func(t *testing.T) {
		defer reset()
		returnErr := errors.New("return error")
		err := NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			return errors.WithStack(returnErr)
		})
		assert.Equal(t, true, rollback)
		assert.Equal(t, false, committed)
		assert.Equal(t, returnErr, errors.Cause(err))
	})
	t.Run("panic error", func(t *testing.T) {
		defer reset()
		panicErr := errors.New("panic error")
		assert.Equal(t, panicErr, NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			panic(panicErr)
		}))
		assert.Equal(t, true, rollback)
		assert.Equal(t, false, committed)
	})

}
