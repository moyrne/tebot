package analyze

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/models"
	"github.com/moyrne/weather"
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
	data, err := weComCn.Get(area)
	if err != nil {
		return "", err
	}
	return "[今日天气]\n" +
		data.Time.Format("2006-01-02") + "\n" +
		data.City + "\n" +
		data.Weather + "\n" +
		data.Wd + " " + data.Ws + "\n" +
		data.Temperature + "℃~" + data.TemperatureN + "℃", nil
}
