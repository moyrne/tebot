package data

import (
	"context"
	"database/sql"

	"github.com/moyrne/tebot/internal/biz"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

// 编译期检验interface实现
var _ biz.UserRepo = userRepo{}

// 持久层实现

type userRepo struct{}

func NewUserRepo() biz.UserRepo {
	return userRepo{}
}

func (u userRepo) GetByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (*biz.User, error) {
	var user biz.User
	query := `select * from user where user_id = ?`
	err := tx.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

func (u userRepo) Save(ctx context.Context, tx dbx.Transaction, user *biz.User) error {
	query := `select * from user where user_id = ? for update`
	err := tx.GetContext(ctx, u, query, user.UserID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	query = `insert into user (user_id,nickname,sex,age) values (?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, user.UserID, user.Nickname, user.Sex, user.Age)
	if err != nil {
		return errors.WithStack(err)
	}
	user.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (u userRepo) UpdateArea(ctx context.Context, tx dbx.Transaction, userID int64, area string) error {
	query := `update user set bind_area = ? where user_id = ?`
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
