package cqhttp

import (
	"context"
	"github.com/moyrne/tebot/internal/biz/cqhttp/template"

	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

var weComCn weather.Weather = weather.WeComCn{}

func SignIn(ctx context.Context, params Params) (string, error) {
	// 记录签到信息
	var area string
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		user, err := data.GetQUserByQUID(ctx, tx, params.QUID)
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
		return (&data.QSignIn{QUID: params.QUID}).Insert(ctx, tx)
	}); err != nil {
		if errors.Is(err, data.ErrAlreadySignIn) {
			return "今日已签到", nil
		}
		return "", err
	}
	// TODO 缓存天气
	data, err := weComCn.Get(area)
	if err != nil {
		return "", err
	}
	return template.Marshal.Template(template.SingInKey).Execute(template.SingInParam(data))
}
