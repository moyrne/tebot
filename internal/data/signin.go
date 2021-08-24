package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.SignInRepo = signInRepo{}

type signInRepo struct{}

func NewSignInRepo() cqhttp.SignInRepo {
	return signInRepo{}
}

func (s signInRepo) GetByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (cqhttp.SignIn, error) {
	var signIn cqhttp.SignIn
	query := `select * from sign_in where user_id = ? and day>= ?`
	err := tx.GetContext(ctx, &signIn, query, userID, time.Now().Format("2006-01-02"))
	if err == nil {
		return signIn, errors.WithStack(cqhttp.ErrAlreadySignIn)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return signIn, errors.WithStack(err)
	}
	return signIn, nil
}

func (s signInRepo) Save(ctx context.Context, tx dbx.Transaction, signIn *cqhttp.SignIn) error {
	if _, err := s.GetByUserID(ctx, tx, signIn.UserID); err != nil {
		return err
	}
	query := `insert into sign_in (user_id,create_at,day) values (?,?,?)`
	result, err := tx.ExecContext(ctx, query, signIn.UserID, time.Now(), time.Now().Format("2006-01-02"))
	if err != nil {
		return errors.WithStack(err)
	}
	signIn.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}
