package models

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// QGroup QQ Group
// 会员才会在表中
type QGroup struct {
	ID   int64  `json:"id"`
	QGID int    `json:"qgid"`
	Name string `json:"name"`
}

func (g QGroup) TableName() string {
	return "q_group"
}

func GetQGroupByQGID(ctx context.Context, tx *sqlx.Tx, qgid sql.NullInt64) (*QGroup, error) {
	var group QGroup
	query := `select * from q_group where qgid = $1`
	if err := tx.GetContext(ctx, &group, query, qgid.Int64); err != nil {
		return nil, errors.WithStack(err)
	}
	return &group, nil
}
