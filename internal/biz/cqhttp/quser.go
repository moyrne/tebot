package cqhttp

import "database/sql"

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
