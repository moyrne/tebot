package data

import (
	"context"
	"database/sql"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"time"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.SignInRepo = SignInData{}
var ErrAlreadySignIn = errors.New("already sign in today")

type SignInData struct{}

func (s SignInData) GetByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (cqhttp.SignIn, error) {
	var signIn cqhttp.SignIn
	query := `select * from q_sign_in where quid = ? and day>= ?`
	err := tx.GetContext(ctx, &signIn, query, userID, time.Now().Format("2006-01-02"))
	if err == nil {
		return signIn, errors.WithStack(ErrAlreadySignIn)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return signIn, errors.WithStack(err)
	}
	return signIn, nil
}

func (s SignInData) Save(ctx context.Context, tx dbx.Transaction, signIn *cqhttp.SignIn) error {
	if _, err := s.GetByUserID(ctx, tx, signIn.UserID); err != nil {
		return err
	}
	query := `insert into q_sign_in (quid,create_at,day) values (?,?,?)`
	result, err := tx.ExecContext(ctx, query, signIn.UserID, time.Now(), time.Now().Format("2006-01-02"))
	if err != nil {
		return errors.WithStack(err)
	}
	signIn.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}
