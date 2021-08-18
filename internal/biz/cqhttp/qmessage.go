package cqhttp

import (
	"context"
	"database/sql"
	"github.com/moyrne/tractor/dbx"
)

// QMessage QQ 聊天记录
type QMessage struct {
	ID          int64         `json:"id"` // 时间戳
	Time        int           `json:"time"`
	SelfID      int           `json:"self_id"`
	PostType    string        `json:"post_type"`    // 上报类型	message
	MessageType string        `json:"message_type"` // 消息类型	private
	SubType     string        `json:"sub_type"`     // 消息子类型 friend,group,group_self,other
	TempSource  string        `json:"temp_source"`  // 临时会话来源
	MessageID   int           `json:"message_id"`   // 消息ID
	GroupID     sql.NullInt64 `json:"group_id"`     // 群ID
	UserID      int           `json:"user_id"`      // 发送者 QQ 号
	Message     string        `json:"message"`      // 消息内容
	RawMessage  string        `json:"raw_message"`  // 原始消息内容
	Font        int           `json:"font"`         // 字体
	Reply       string        `json:"reply"`        // 回复

	QUser *QUser `json:"q_user"` // User 信息
}

type QMessageRepo interface {
	Save(ctx context.Context, tx dbx.Transaction, message *QMessage) error
	SetReply(ctx context.Context, tx dbx.Transaction, id int64, reply string) error
}

func NewQMessageUseCase(repo QMessageRepo) *QMessageUseCase {
	return &QMessageUseCase{repo: repo}
}

type QMessageUseCase struct {
	repo QMessageRepo
}
