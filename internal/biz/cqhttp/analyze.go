package cqhttp

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
)

type QMessageRepo interface {
	Event(message *QMessage) (string, error)
}

func NewQMessageUsecase(repo QMessageRepo) *QMessageUsecase {
	return &QMessageUsecase{repo: repo}
}

type QMessageUsecase struct {
	repo QMessageRepo
}

func (uc *QMessageUsecase) Event(message *QMessage) (string, error) {
	return uc.repo.Event(message)
}

// QMessage QQ 聊天记录
type QMessage struct {
	ID          int64         `json:"id"` // 时间戳
	Time        int           `json:"time"`
	SelfID      int           `json:"self_id"`
	PostType    string        `json:"post_type"`    // 上报类型	message
	MessageType string        `json:"message_type"` // 消息类型	private
	SubType     string        `json:"sub_type"`     // 消息子类型 friend,group,group_self,other
	TempSource  string        `json:"temp_source"`  // 临时会话来源
	MessageID   int           `json:"message_id"`   // 消息ID
	GroupID     sql.NullInt64 `json:"group_id"`     // 群ID
	UserID      int           `json:"user_id"`      // 发送者 QQ 号
	Message     string        `json:"message"`      // 消息内容
	RawMessage  string        `json:"raw_message"`  // 原始消息内容
	Font        int           `json:"font"`         // 字体
	Reply       string        `json:"reply"`        // 回复

	//QUser *QUser `json:"q_user"` // User 信息
}

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
