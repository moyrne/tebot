package data

import (
	"context"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

type Reply struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Msg    string `json:"msg"` // 前缀搜索
	Reply  string `json:"reply"`
}

func (r Reply) TableName() string {
	return "q_reply"
}

func SelectQReply(ctx context.Context, tx dbx.Transaction) (replies []Reply, err error) {
	query := `select * from q_reply`
	err = tx.SelectContext(ctx, &replies, query)
	return replies, errors.WithStack(err)
}
