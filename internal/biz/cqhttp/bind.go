package cqhttp

import (
	"context"
	"strings"

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
		return uc.user.UpdateArea(ctx, tx, m.UserID, area)
	}); err != nil {
		return "", err
	}
	return "绑定成功：" + area, nil
}
