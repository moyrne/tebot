package data

import (
	"context"

	"github.com/moyrne/tebot/internal/biz"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ biz.GroupRepo = groupRepo{}

type groupRepo struct{}

func NewGroupRepo() biz.GroupRepo {
	return groupRepo{}
}

func (g groupRepo) GetByID(ctx context.Context, tx dbx.Transaction, groupID int64) (*biz.Group, error) {
	var group biz.Group
	query := `select * from group where group_id = ?`
	if err := tx.GetContext(ctx, &group, query, groupID); err != nil {
		return nil, errors.WithStack(err)
	}
	return &group, nil
}
