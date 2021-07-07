package models

// QUser QQ User
type QUser struct {
	ID       int    `json:"id"`
	QUID     int    `json:"quid"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
}
