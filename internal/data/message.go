package data

import (
	"context"

	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.MessageRepo = messageRepo{}

// messageRepo 聊天记录持久化
type messageRepo struct{}

func NewMessageRepo() cqhttp.MessageRepo {
	return messageRepo{}
}

func (q messageRepo) Save(ctx context.Context, tx dbx.Transaction, m *cqhttp.Message) error {
	query := `insert into message (time,self_id,post_type,message_type,sub_type,temp_source,message_id,group_id,user_id,message,raw_message,font) values (?,?,?,?,?,?,?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, m.Time, m.SelfID, m.PostType, m.MessageType, m.SubType, m.TempSource, m.MessageID, m.GroupID, m.UserID, m.Message, m.RawMessage, m.Font)
	if err != nil {
		return errors.WithStack(err)
	}
	m.ID, err = result.LastInsertId()
	return errors.WithStack(err)
}

func (q messageRepo) SetReply(ctx context.Context, tx dbx.Transaction, id int64, reply string) error {
	query := `update q_message set reply = ? where id = ?`
	r, err := tx.ExecContext(ctx, query, reply, id)
	if err != nil {
		return errors.WithStack(err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	if affected == 0 {
		return errors.WithStack(errors.Wrapf(database.ErrRowsAffectedZero, "id: %d,reply: %s", id, reply))
	}
	return nil
}
