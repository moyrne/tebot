package data

import (
	"context"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

var _ cqhttp.GroupRepo = groupRepo{}

type groupRepo struct{}

func NewGroupRepo() cqhttp.GroupRepo {
	return groupRepo{}
}

func (g groupRepo) GetByID(ctx context.Context, tx dbx.Transaction, groupID int64) (*cqhttp.Group, error) {
	var group cqhttp.Group
	query := `select * from group where group_id = ?`
	if err := tx.GetContext(ctx, &group, query, groupID); err != nil {
		return nil, errors.WithStack(err)
	}
	return &group, nil
}
