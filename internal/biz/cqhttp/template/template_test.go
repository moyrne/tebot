package template

import (
	"testing"
	"time"

	"github.com/moyrne/weather"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	ti, err := time.Parse("2006-01-02 15:04:05", "2021-07-11 14:13:00")
	assert.Equal(t, nil, err)
	reply, err := Marshal.Template(SingInKey).Execute(SingInParam(weather.Data{
		City:         "深圳",
		Temperature:  "27",
		TemperatureN: "31",
		Weather:      "多云",
		Wd:           "南风",
		Ws:           "3级",
		Time:         ti,
	}))
	assert.Equal(t, nil, err)
	assert.Equal(t, "签到成功！\n๑ 今日天气\n๑ 2021-07-11\n๑ 深圳\n๑ 多云\n๑ 南风  3级\n๑ 27℃~31℃", reply)
}
