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

func (uc *QMessageUseCase) Event(ctx context.Context, m *QMessage) (string, error) {
	err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		return uc.repo.Save(ctx, tx, m)
	})
	if err != nil {
		logs.Error("insert CqHTTP params failed", "error", err)
	}

	// User Filter, 检查是否被禁用
	if m.QUser.Ban {
		logs.Info("quser has been band", "user", m.QUser)
		return "", nil
	}
	switch m.MessageType {
	// TODO const
	case "MTPrivate":
		m.Reply, err = uc.private(ctx, m)
	case "MTGroup":
		m.Reply, err = uc.group(ctx, m)
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
		return uc.repo.SetReply(ctx, tx, m.ID, m.Reply)
	}); err != nil {
		logs.Error("update cqhttp params failed", "error", err)
	}

	return m.Reply, nil
}

func (uc *QMessageUseCase) private(ctx context.Context, m *QMessage) (string, error) {
	return uc.doEvent(ctx, m)
}

func (uc *QMessageUseCase) group(ctx context.Context, m *QMessage) (string, error) {
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		//_, err := data.GetQGroupByQGID(ctx, tx, params.GroupID)
		//return err
		return nil
	}); err != nil {
		// 只显示主要原因, 抛弃栈信息
		return "", errors.WithMessage(errors.Cause(err), strconv.Itoa(int(m.GroupID.Int64)))
	}
	return uc.doEvent(ctx, m)
	// TODO at Sender
	//return api.Reply{Reply: reply, ATSender: true}, nil
}

func (uc *QMessageUseCase) doEvent(ctx context.Context, m *QMessage) (string, error) {
	resp, err := rangeDo(ctx, Functions, m)
	if err == nil {
		return resp, nil
	}
	if !errors.Is(err, ErrNotMatch) {
		return "", err
	}

	// 等待3秒 回复
	defer time.Sleep(time.Second * 3)
	defer logs.Info("reply", "content", m.Reply)

	// TODO 匹配简单回复
	return SimpleReply(ctx, m)
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
	Name  string
	Fn    func(context.Context, *QMessage) (string, error)
	Match func(s, v string) bool
}

var (
	Functions = []Menu{
		{Name: "menu", Fn: PrintMenu, Match: Equal},
		{Name: "签到", Fn: SignIn, Match: Equal},
		{Name: "绑定位置", Fn: BindArea, Match: Prefix},
	}
)

var ErrNotMatch = errors.New("not match")

func rangeDo(ctx context.Context, functions []Menu, m *QMessage) (string, error) {
	for _, menu := range functions {
		if menu.Match(m.Message, menu.Name) {
			if err := rateLimiter.Rate(ctx, menu.Name, m.UserID); err != nil {
				return "", err
			}
			return menu.Fn(ctx, m)
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
