package biz

import (
	"context"

	"github.com/moyrne/tractor/dbx"
)

type Group struct {
	ID      int64  `json:"id"`
	GroupID int64  `json:"group_id" db:"group_id"`
	Name    string `json:"name"`
}

type GroupRepo interface {
	GetByID(ctx context.Context, tx dbx.Transaction, groupID int64) (*Group, error)
}
