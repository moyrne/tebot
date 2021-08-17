package cqhttp

import (
	"context"
	"strings"

	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tractor/dbx"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

var replacer = strings.NewReplacer("绑定位置", "", " ", "", "\t", "")

func BindArea(ctx context.Context, params Params) (string, error) {
	area := replacer.Replace(params.Message)
	if area == "" {
		return "模板：绑定地区 深圳", nil
	}
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		_, err := weather.GetCityID(area)
		if err != nil {
			return errors.WithStack(err)
		}
		return data.UpdateArea(ctx, tx, params.QUID, area)
	}); err != nil {
		return "", err
	}
	return "绑定成功：" + area, nil
}
