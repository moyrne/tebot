package analyze

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/models"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

var weComCn weather.Weather = weather.WeComCn{}

func SignIn(ctx context.Context, params Params) (string, error) {
	// 记录签到信息
	var area string
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		user, err := models.GetQUserByQUID(ctx, tx, params.QUID)
		if err != nil {
			return err
		}
		area = user.BindArea.String
		return nil
	}); err != nil {
		return "", err
	}
	if area == "" {
		return "未绑定位置\n" +
			"例如: 绑定位置 深圳", nil
	}
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return (&models.QSignIn{QUID: params.QUID}).Insert(ctx, tx)
	}); err != nil {
		if errors.Is(err, models.ErrAlreadySignIn) {
			return "今日已签到", nil
		}
		return "", err
	}
	// TODO 缓存天气
	data, err := weComCn.Get(area)
	if err != nil {
		return "", err
	}
	return "签到成功！\n" +
		"[今日天气]\n" +
		data.Time.Format("2006-01-02") + "\n" +
		data.City + "\n" +
		data.Weather + "\n" +
		data.Wd + " " + data.Ws + "\n" +
		data.Temperature + "℃~" + data.TemperatureN + "℃", nil
}
