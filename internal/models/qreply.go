package models

import (
	"context"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

type QReply struct {
	ID    int64  `json:"id"`
	QUID  int    `json:"quid"`
	Msg   string `json:"msg"` // 前缀搜索
	Reply string `json:"reply"`
}

func (r QReply) TableName() string {
	return "q_reply"
}

func SelectQReply(ctx context.Context, tx dbx.Transaction) (replies []QReply, err error) {
	query := `select * from q_reply`
	err = tx.SelectContext(ctx, &replies, query)
	return replies, errors.WithStack(err)
}
