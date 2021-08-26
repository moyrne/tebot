package biz

import (
	"context"
	"time"

	"github.com/moyrne/tractor/dbx"
)

type SignIn struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Day      string    `json:"day"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}

type SignInRepo interface {
	GetByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (SignIn, error)
	Save(ctx context.Context, tx dbx.Transaction, signIn *SignIn) error
}
