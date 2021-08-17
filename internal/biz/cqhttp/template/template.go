package template

import (
	"bytes"
	"text/template"

	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

const (
	SingInKey = `signIn`
)

var Marshal ReplyMarshal

func init() {
	Marshal = ReplyMarshal{temps: map[string]*Template{}}
	Marshal.Register("signIn", `签到成功！
๑ 今日天气
๑ {{.Time}}
๑ {{.City}}
๑ {{.Weather}}
๑ {{.Wd}}  {{.Ws}}
๑ {{.Temperature}}℃~{{.TemperatureN}}℃`)
}

type ReplyMarshal struct {
	temps map[string]*Template
}

func (m *ReplyMarshal) Register(name, temp string) {
	marshal, err := template.New(name).Parse(temp)
	if err != nil {
		logs.Panic("parse template", "error", err)
	}
	m.temps[name] = &Template{Template: marshal}
}

func (m *ReplyMarshal) Template(name string) *Template {
	return m.temps[name]
}

type Template struct {
	*template.Template
}

func (t Template) Execute(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}

type SignInData struct {
	City         string `json:"city"`          // 城市
	Temperature  string `json:"temperature"`   // 最低气温
	TemperatureN string `json:"temperature_n"` // 最高气温
	Weather      string `json:"weather"`       // 天气
	Wd           string `json:"wd"`            // 风向
	Ws           string `json:"ws"`            // 风速
	Time         string `json:"time"`          // 时间
}

func SingInParam(data weather.Data) SignInData {
	return SignInData{
		City:         data.City,
		Temperature:  data.Temperature,
		TemperatureN: data.TemperatureN,
		Weather:      data.Weather,
		Wd:           data.Wd,
		Ws:           data.Ws,
		Time:         data.Time.Format("2006-01-02"),
	}
}
