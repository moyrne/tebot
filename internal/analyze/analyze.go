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
	Name string
	Fn   func(context.Context, Params) (string, error)
}

var (
	EqualFunctions = []Menu{
		{Name: "签到", Fn: SignIn},
	}
	PrefixFunctions = []Menu{
		{Name: "绑定位置", Fn: BindArea},
	}
)

func Analyze(ctx context.Context, params Params) (string, error) {
	// TODO 优先匹配高级功能
	// 1. 相等
	resp, err := rangeDo(ctx, EqualFunctions, params, equal)
	if err == nil {
		return resp, nil
	}
	if !errors.Is(err, ErrNotMatch) {
		return "", err
	}
	// TODO 2. 包含
	resp, err = rangeDo(ctx, PrefixFunctions, params, func(msg, name string) bool {
		return strings.HasPrefix(strings.ReplaceAll(msg, " ", ""), name)
	})
	if err == nil {
		return resp, nil
	}
	if !errors.Is(err, ErrNotMatch) {
		return "", err
	}

	// TODO 2. 正则

	// TODO 匹配简单回复

	// 未匹配到
	return "", errors.WithStack(ErrNotMatch)
}

func equal(s, v string) bool {
	return strings.ReplaceAll(s, " ", "") == strings.ReplaceAll(v, " ", "")
}

var ErrNotMatch = errors.New("not match")

func rangeDo(ctx context.Context, functions []Menu, params Params, match func(msg, name string) bool) (string, error) {
	for _, menu := range functions {
		if match(params.Message, menu.Name) {
			if err := rateLimiter.Rate(ctx, menu.Name, params.QUID); err != nil {
				return "", err
			}
			return menu.Fn(ctx, params)
		}
	}
	return "", errors.WithStack(ErrNotMatch)
}
