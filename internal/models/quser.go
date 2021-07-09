package models

// QUser QQ User
type QUser struct {
	ID       int64  `json:"id"`
	QUID     int    `json:"quid"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`

	BindArea string `json:"bind_area"` // 所在地
	Mode     string `json:"mode"`      // 人设模式
}

func (u QUser) TableName() string {
	return "q_user"
}
