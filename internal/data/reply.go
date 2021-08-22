package data

import (
	"context"
	"github.com/moyrne/tebot/internal/biz/cqhttp"

	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.ReplyRepo = replyRepo{}

type replyRepo struct{}

func NewReplyRepo() cqhttp.ReplyRepo {
	return replyRepo{}
}

func (r replyRepo) Replies(ctx context.Context, tx dbx.Transaction) (replies []cqhttp.Reply, err error) {
	query := `select * from reply`
	err = tx.SelectContext(ctx, &replies, query)
	return replies, errors.WithStack(err)
}
