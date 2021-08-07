package models

import (
	"context"
	"database/sql"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

// QUser QQ User
type QUser struct {
	ID       int64  `json:"id"`
	QUID     int    `json:"quid"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`

	BindArea sql.NullString `json:"bind_area" db:"bind_area"` // 所在地
	Mode     sql.NullString `json:"mode"`                     // 人设模式

	Ban bool `json:"ban"` // 被禁
}

func (u QUser) TableName() string {
	return "q_user"
}

func GetQUserByQUID(ctx context.Context, tx dbx.Transaction, quid int) (*QUser, error) {
	var user QUser
	query := `select * from q_user where quid = ?`
	err := tx.GetContext(ctx, &user, query, quid)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

func (u *QUser) GetOrInsert(ctx context.Context, tx dbx.Transaction) error {
	query := `select * from q_user where quid = ? for update`
	err := tx.GetContext(ctx, u, query, u.QUID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	query = `insert into q_user (quid,nickname,sex,age) values (?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, u.QUID, u.Nickname, u.Sex, u.Age)
	if err != nil {
		return errors.WithStack(err)
	}
	u.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func UpdateArea(ctx context.Context, tx dbx.Transaction, quid int, area string) error {
	query := `update q_user set bind_area = ? where quid = ?`
	result, err := tx.ExecContext(ctx, query, area, quid)
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
