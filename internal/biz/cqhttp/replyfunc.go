package cqhttp

import (
	"context"
	"strings"

	"github.com/moyrne/tebot/internal/biz/cqhttp/template"
	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

var replacer = strings.NewReplacer("绑定位置", "", " ", "", "\t", "")

func BindAreaMethod(ctx context.Context, uc *EventUseCase, m *Message) (string, error) {
	area := replacer.Replace(m.Message)
	if area == "" {
		return "模板：绑定地区 深圳", nil
	}
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		_, err := weather.GetCityID(area)
		if err != nil {
			return errors.WithStack(err)
		}
		return uc.repo.UpdateUserArea(ctx, tx, m.UserID, area)
	}); err != nil {
		return "", err
	}
	return "绑定成功：" + area, nil
}

var weComCn weather.Weather = weather.WeComCn{}

func SignInMethod(ctx context.Context, uc *EventUseCase, m *Message) (string, error) {
	// 记录签到信息
	var area string
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		user, err := uc.repo.GetUserByUserID(ctx, tx, m.UserID)
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
		return uc.repo.SaveSignIn(ctx, tx, &SignIn{UserID: m.UserID})
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

func PrintMenu(_ context.Context, _ *EventUseCase, _ *Message) (string, error) {
	return "๑ 菜单\n" +
		"๑ 1.绑定位置;(绑定位置 深圳)\n" +
		"๑ 2.签到", nil
}
