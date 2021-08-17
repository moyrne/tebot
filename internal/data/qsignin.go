package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
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

func (s *QSignIn) Insert(ctx context.Context, tx dbx.Transaction) error {
	if err := s.GetQSignInByQUID(ctx, tx); err != nil {
		return err
	}
	query := `insert into q_sign_in (quid,create_at,day) values (?,?,?)`
	result, err := tx.ExecContext(ctx, query, s.QUID, time.Now(), time.Now().Format("2006-01-02"))
	if err != nil {
		return errors.WithStack(err)
	}
	s.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (s *QSignIn) GetQSignInByQUID(ctx context.Context, tx dbx.Transaction) error {
	query := `select * from q_sign_in where quid = ? and day>= ?`
	err := tx.GetContext(ctx, s, query, s.QUID, time.Now().Format("2006-01-02"))
	if err == nil {
		return errors.WithStack(ErrAlreadySignIn)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	return nil
}
