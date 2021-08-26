package biz

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/pkg/autoreply"
	"github.com/moyrne/tebot/internal/pkg/template"
	"github.com/moyrne/tractor/dbx"
	"github.com/moyrne/weather"
	"github.com/pkg/errors"
)

var ErrAlreadySignIn = errors.New("already sign in today")

var replacer = strings.NewReplacer("绑定位置", "", " ", "", "\t", "")

func BindAreaMethod(uc *EventUseCase) func(ctx context.Context, m *autoreply.Message) (string, error) {
	return func(ctx context.Context, m *autoreply.Message) (string, error) {
		userID, err := strconv.Atoi(m.UserID)
		if err != nil {
			return "", errors.WithStack(err)
		}
		area := replacer.Replace(m.Message)
		if area == "" {
			return "模板：绑定地区 深圳", nil
		}
		if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
			_, err := weather.GetCityID(area)
			if err != nil {
				return errors.WithStack(err)
			}
			return uc.repo.UpdateUserArea(ctx, tx, int64(userID), area)
		}); err != nil {
			return "", err
		}
		return "绑定成功：" + area, nil
	}
}

var (
	weComCn    weather.Weather = weather.WeComCn{}
	signInTemp *template.Template
)

func init() {
	var err error
	signInTemp, err = template.NewTemplate("SignIn", `签到成功！
๑ 今日天气
๑ {{.Time}}
๑ {{.City}}
๑ {{.Weather}}
๑ {{.Wd}}  {{.Ws}}
๑ {{.Temperature}}℃~{{.TemperatureN}}℃`)
	if err != nil {
		log.Panicf("parse signin template error %v", errors.WithStack(err))
	}
}

func SignInMethod(uc *EventUseCase) func(ctx context.Context, m *autoreply.Message) (string, error) {
	return func(ctx context.Context, m *autoreply.Message) (string, error) {
		userID, err := strconv.Atoi(m.UserID)
		if err != nil {
			return "", errors.WithStack(err)
		}
		// 记录签到信息
		var area string
		if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
			user, err := uc.repo.GetUserByUserID(ctx, tx, int64(userID))
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
			return uc.repo.SaveSignIn(ctx, tx, &SignIn{UserID: int64(userID)})
		}); err != nil {
			if errors.Is(err, ErrAlreadySignIn) {
				return "今日已签到", nil
			}
			return "", err
		}
		// TODO 缓存天气
		wt, err := weComCn.Get(area)
		if err != nil {
			return "", err
		}
		return signInTemp.Execute(signParam(wt))
	}
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

func signParam(data weather.Data) SignInData {
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

func PrintMenu(_ context.Context, _ *autoreply.Message) (string, error) {
	return "๑ 菜单\n" +
		"๑ 1.绑定位置;(绑定位置 深圳)\n" +
		"๑ 2.签到", nil
}
