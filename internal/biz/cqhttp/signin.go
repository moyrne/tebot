package cqhttp

import (
	"context"
	"time"

	"github.com/moyrne/tebot/internal/biz/cqhttp/template"
	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
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

var weComCn weather.Weather = weather.WeComCn{}

func SignInMethod(ctx context.Context, uc *EventUseCase, m *Message) (string, error) {
	// 记录签到信息
	var area string
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		user, err := uc.user.GetByUserID(ctx, tx, m.UserID)
		if err != nil {
			return err
		}
		area = user.BindArea.String
		return nil
	}); err != nil {
		return "", err
	}
	if area == "" {
		return "未绑定位置\n例如: 绑定位置 深圳", nil
	}
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		return uc.signIn.Save(ctx, tx, &SignIn{UserID: m.UserID})
	}); err != nil {
		if errors.Is(err, data.ErrAlreadySignIn) {
			return "今日已签到", nil
		}
		return "", err
	}
	// TODO 缓存天气
	wt, err := weComCn.Get(area)
	if err != nil {
		return "", err
	}
	return template.Marshal.Template(template.SingInKey).Execute(template.SingInParam(wt))
}
