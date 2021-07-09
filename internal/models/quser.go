package models

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// QUser QQ User
type QUser struct {
	ID       int64  `json:"id"`
	QUID     int    `json:"quid"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`

	BindArea string `json:"bind_area"` // 所在地
	Mode     string `json:"mode"`      // 人设模式

	Ban bool `json:"ban"` // 被禁
}

func (u QUser) TableName() string {
	return "q_user"
}

func (u *QUser) GetOrInsert(ctx context.Context, tx *sqlx.Tx) error {
	query := `select * from q_user where quid = $1 for update`
	err := tx.GetContext(ctx, u, query, u.QUID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.WithStack(err)
	}
	query = `insert into q_user (quid,nickname,sex,age) values ($1,$2,$3,$4) returning id`
	return errors.WithStack(tx.GetContext(ctx, &u.ID, query, u.QUID, u.Nickname, u.Sex, u.Age))
}
