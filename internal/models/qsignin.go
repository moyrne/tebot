package models

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type QSignIn struct {
	ID       int64     `json:"id"`
	QUID     int       `json:"quid"`
	Day      string    `json:"day"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}

func (s QSignIn) TableName() string {
	return "q_sign_in"
}

var ErrAlreadySignIn = errors.New("already sign in today")

func (s *QSignIn) Insert(ctx context.Context, tx *sqlx.Tx) error {
	if err := s.GetQSignInByQUID(ctx, tx); err != nil {
		return err
	}
	query := `insert into q_sign_in (quid,create_at,day) values ($1,$2,$3) returning id`
	if err := tx.GetContext(ctx, &s.ID, query, s.QUID, time.Now(), time.Now().Format("2006-01-02")); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *QSignIn) GetQSignInByQUID(ctx context.Context, tx *sqlx.Tx) error {
	query := `select * from q_sign_in where quid = $1 and create_at > $2`
	err := tx.GetContext(ctx, s, query, s.QUID, time.Now())
	if err == nil {
		return errors.WithStack(ErrAlreadySignIn)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	return nil
}
