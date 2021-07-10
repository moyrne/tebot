package analyze

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/models"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
	"strings"
)

var replacer = strings.NewReplacer("绑定位置", "", " ", "", "\t", "")

func BindArea(ctx context.Context, params Params) (string, error) {
	area := replacer.Replace(params.Message)
	if area == "" {
		return "模板：绑定地区 深圳", nil
	}
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := weather.GetCityID(area)
		if err != nil {
			return errors.WithStack(err)
		}
		return models.UpdateArea(ctx, tx, params.QUID, area)
	}); err != nil {
		return "", err
	}
	return "绑定成功：" + area, nil
}
