package data

import (
	"context"
	"database/sql"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

// 编译期检验interface实现
var _ cqhttp.UserRepo = UserData{}

// 持久层实现

type UserData struct{}

func (u UserData) GetByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (*cqhttp.User, error) {
	var user cqhttp.User
	query := `select * from q_user where quid = ?`
	err := tx.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

func (u UserData) Save(ctx context.Context, tx dbx.Transaction, user *cqhttp.User) error {
	query := `select * from q_user where quid = ? for update`
	err := tx.GetContext(ctx, u, query, user.UserID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	query = `insert into q_user (quid,nickname,sex,age) values (?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, user.UserID, user.Nickname, user.Sex, user.Age)
	if err != nil {
		return errors.WithStack(err)
	}
	user.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (u UserData) UpdateArea(ctx context.Context, tx dbx.Transaction, userID int64, area string) error {
	query := `update q_user set bind_area = ? where quid = ?`
	result, err := tx.ExecContext(ctx, query, area, userID)
	if err != nil {
		return errors.WithStack(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	if affected == 0 {
		return errors.WithStack(database.ErrRowsAffectedZero)
	}
	return nil
}
