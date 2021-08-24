package cqhttp

import (
	"context"
	"database/sql"

	"github.com/moyrne/tractor/dbx"
)

// User QQ User
type User struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`

	BindArea sql.NullString `json:"bind_area" db:"bind_area"` // 所在地
	Mode     sql.NullString `json:"mode"`                     // 人设模式

	Ban bool `json:"ban"` // 被禁
}

type UserRepo interface {
	GetByUserID(ctx context.Context, tx dbx.Transaction, id int64) (*User, error)
	Save(ctx context.Context, tx dbx.Transaction, u *User) error
	UpdateArea(ctx context.Context, tx dbx.Transaction, userID int64, area string) error
}
