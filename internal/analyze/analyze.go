package analyze

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

type Params struct {
	QUID    int    `json:"quid"`
	Message string `json:"message"`
}

type Menu struct {
	Name  string
	Fn    func(context.Context, Params) (string, error)
	Match func(s, v string) bool
}

var (
	Functions = []Menu{
		{Name: "menu", Fn: PrintMenu, Match: Equal},
		{Name: "签到", Fn: SignIn, Match: Equal},
		{Name: "绑定位置", Fn: BindArea, Match: Prefix},
	}
)

func Analyze(ctx context.Context, params Params) (string, error) {
	// TODO 优先匹配高级功能
	// 相等
	resp, err := rangeDo(ctx, Functions, params)
	if err == nil {
		return resp, nil
	}
	if !errors.Is(err, ErrNotMatch) {
		return "", err
	}

	// TODO 匹配简单回复
	return SimpleReply(ctx, params)
}

var ErrNotMatch = errors.New("not match")

func rangeDo(ctx context.Context, functions []Menu, params Params) (string, error) {
	for _, menu := range functions {
		if menu.Match(params.Message, menu.Name) {
			if err := rateLimiter.Rate(ctx, menu.Name, params.QUID); err != nil {
				return "", err
			}
			return menu.Fn(ctx, params)
		}
	}
	return "", errors.WithStack(ErrNotMatch)
}

func Equal(s, v string) bool {
	return strings.ReplaceAll(s, " ", "") == strings.ReplaceAll(v, " ", "")
}

func Prefix(s, v string) bool {
	return strings.HasPrefix(strings.ReplaceAll(s, " ", ""), v)
}
