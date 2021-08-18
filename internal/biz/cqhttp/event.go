package cqhttp

import (
	"context"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tractor/dbx"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TODO log error set to return error

func (uc *EventUseCase) Event(ctx context.Context, m *Message) (string, error) {
	err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		// TODO get or create user
		if err := uc.user.Save(ctx, tx, m.User); err != nil {
			return err
		}
		return uc.message.Save(ctx, tx, m)
	})
	if err != nil {
		logs.Error("insert CqHTTP params failed", "error", err)
	}

	// User Filter, 检查是否被禁用
	if m.User.Ban {
		logs.Info("quser has been band", "user", m.User)
		return "", nil
	}
	switch m.MessageType {
	// TODO const
	case "MTPrivate":
		m.Reply, err = uc.privateHandler(ctx, m)
	case "MTGroup":
		m.Reply, err = uc.groupHandler(ctx, m)
	default:
		// log error
		logs.Error("unsupported message_type", "type", m.MessageType)
		return "", nil
	}

	if errors.Is(err, ErrNotMatch) {
		return "", nil
	}
	if err != nil {
		logs.Info("get reply", "error", err)
		return "", nil
	}

	if err = database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		return uc.message.SetReply(ctx, tx, m.ID, m.Reply)
	}); err != nil {
		logs.Error("update cqhttp params failed", "error", err)
	}

	return m.Reply, nil
}

func (uc *EventUseCase) privateHandler(ctx context.Context, m *Message) (string, error) {
	return uc.doEvent(ctx, m)
}

func (uc *EventUseCase) groupHandler(ctx context.Context, m *Message) (string, error) {
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		_, err := uc.group.GetByID(ctx, tx, m.GroupID.Int64)
		return err
	}); err != nil {
		// 只显示主要原因, 抛弃栈信息
		return "", errors.WithMessage(errors.Cause(err), strconv.Itoa(int(m.GroupID.Int64)))
	}
	return uc.doEvent(ctx, m)
	// TODO at Sender
	//return api.Replies{Replies: reply, ATSender: true}, nil
}

func (uc *EventUseCase) doEvent(ctx context.Context, m *Message) (string, error) {
	// 等待3秒 回复
	defer time.Sleep(time.Second * 3)
	defer logs.Info("reply", "content", m.Reply)

	// TODO 匹配简单回复
	return ReplyMethod(ctx, uc, m)
}

const (
	TSGroup       = "group"       // 群聊
	TSConsult     = "consult"     // QQ咨询
	TSSearch      = "search"      // 查找
	TSFilm        = "film"        // QQ电影
	TSHotTalk     = "hottalk"     // 热聊
	TSVerify      = "verify"      // 验证消息
	TSMultiChat   = "multichat"   // 多人聊天
	TSAppointment = "appointment" // 约会
	TSMailList    = "maillist"    // 通讯录
)

type Menu struct {
	Name string
	Fn   func(ctx context.Context, uc *EventUseCase, m *Message) (string, error)
}

type Match struct {
	Name string
	Fn   func(s, v string) bool
}

// TODO 抽离 数据库内容，封装到 pkg 中
var (
	Functions = map[string]Menu{
		"PrintMenu":      {Name: "menu", Fn: PrintMenu},
		"SignInMethod":   {Name: "签到", Fn: SignInMethod},
		"BindAreaMethod": {Name: "绑定位置", Fn: BindAreaMethod},
	}
	Matches = map[string]Match{
		"Equal":  {Name: "Equal", Fn: Equal},
		"Prefix": {Name: "Prefix", Fn: Prefix},
	}
)

var ErrNotMatch = errors.New("not match")

func Equal(s, v string) bool {
	return strings.ReplaceAll(s, " ", "") == strings.ReplaceAll(v, " ", "")
}

func Prefix(s, v string) bool {
	return strings.HasPrefix(strings.ReplaceAll(s, " ", ""), v)
}
