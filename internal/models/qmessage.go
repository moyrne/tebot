package models

import (
	"context"
	"database/sql"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

const (
	TSGroup       = "group"       // 群聊
	TSConsult     = "consult"     // QQ咨询
	TSSearch      = "search"      // 查找
	TSFilm        = "film"        // QQ电影
	TSHotTalk     = "hottalk"     // 热聊
	TSVerify      = "verify"      // 验证消息
	TSMultiChat   = "multichat"   // 多人聊天
	TSAppointment = "appointment" // 约会
	TSMailList    = "maillist"    // 通讯录
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

func (m QMessage) TableName() string {
	return "q_message"
}

func (m *QMessage) Insert(ctx context.Context, tx dbx.Transaction) error {
	if err := m.QUser.GetOrInsert(ctx, tx); err != nil {
		return err
	}
	query := `insert into q_message (time,self_id,post_type,message_type,sub_type,temp_source,message_id,group_id,user_id,message,raw_message,font) values (?,?,?,?,?,?,?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, m.Time, m.SelfID, m.PostType, m.MessageType, m.SubType, m.TempSource, m.MessageID, m.GroupID, m.UserID, m.Message, m.RawMessage, m.Font)
	if err != nil {
		return errors.WithStack(err)
	}
	m.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (m *QMessage) SetReply(ctx context.Context, tx dbx.Transaction) error {
	query := `update q_message set reply = ? where id = ?`
	r, err := tx.ExecContext(ctx, query, m.Reply, m.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	if affected == 0 {
		return errors.WithStack(errors.Wrapf(database.ErrRowsAffectedZero, "m: %#v", *m))
	}
	return nil
}
