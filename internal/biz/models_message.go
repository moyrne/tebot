package biz

import (
	"context"
	"database/sql"

	"github.com/moyrne/tractor/dbx"
)

// Message QQ 聊天记录
type Message struct {
	ID          int64         `json:"id"` // 时间戳
	Time        int           `json:"time"`
	SelfID      int           `json:"self_id" db:"self_id"`
	PostType    string        `json:"post_type" db:"post_type"`       // 上报类型	message
	MessageType string        `json:"message_type" db:"message_type"` // 消息类型	private
	SubType     string        `json:"sub_type" db:"sub_type"`         // 消息子类型 friend,group,group_self,other
	TempSource  string        `json:"temp_source" db:"temp_source"`   // 临时会话来源
	MessageID   int           `json:"message_id" db:"message_id"`     // 消息ID
	GroupID     sql.NullInt64 `json:"group_id" db:"group_id"`         // 群ID
	UserID      int64         `json:"user_id" db:"user_id"`           // 发送者 QQ 号
	Message     string        `json:"message"`                        // 消息内容
	RawMessage  string        `json:"raw_message" db:"raw_message"`   // 原始消息内容
	Font        int           `json:"font"`                           // 字体
	Reply       string        `json:"reply"`                          // 回复

	User *User `json:"user"` // User 信息
}

type MessageRepo interface {
	Save(ctx context.Context, tx dbx.Transaction, message *Message) error
	SetReply(ctx context.Context, tx dbx.Transaction, id int64, reply string) error
}
