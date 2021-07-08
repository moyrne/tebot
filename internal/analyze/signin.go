package analyze

import (
	"context"
	"github.com/moyrne/weather"
)

var weComCn weather.Weather = weather.WeComCn{}

func SignIn(ctx context.Context, params Params) (string, error) {
	// 记录签到信息
	dfPosition := "深圳"
	data, err := weComCn.Get(dfPosition)
	if err != nil {
		return "", err
	}
	return "[今日天气]\n" +
		data.Time.Format("2006-01-02 15:04:05") + "\n" +
		data.City + "\n" +
		data.Weather + "\n" +
		data.Wd + " " + data.Ws + "\n" +
		data.Temperature + "℃~" + data.TemperatureN + "℃", nil
}
