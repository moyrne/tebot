package biz

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/pkg/autoreply"
	"github.com/moyrne/tebot/internal/pkg/ratelimit"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EventRepo interface {
	SaveMessage(ctx context.Context, tx dbx.Transaction, message *Message) error
	SetMessageReply(ctx context.Context, tx dbx.Transaction, id int64, reply string) error
	GetGroupByID(ctx context.Context, tx dbx.Transaction, groupID int64) (*Group, error)
	GetSignInByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (SignIn, error)
	SaveSignIn(ctx context.Context, tx dbx.Transaction, signIn *SignIn) error
	GetUserByUserID(ctx context.Context, tx dbx.Transaction, id int64) (*User, error)
	SaveUser(ctx context.Context, tx dbx.Transaction, u *User) error
	UpdateUserArea(ctx context.Context, tx dbx.Transaction, userID int64, area string) error
	Log(ctx context.Context, tx dbx.Transaction, log *Log) error
}

func NewEventUseCase(repo EventRepo) *EventUseCase {
	useCase := &EventUseCase{repo: repo}
	autoreply.RegisterMatches("Equal", Equal)
	autoreply.RegisterMatches("Prefix", Prefix)

	autoreply.RegisterFunctions("PrintMenu", PrintMenu)
	autoreply.RegisterFunctions("SignInMethod", SignInMethod(useCase))
	autoreply.RegisterFunctions("BindAreaMethod", BindAreaMethod(useCase))
	return useCase
}

type EventUseCase struct {
	repo EventRepo
}

var (
	ErrUserBand           = errors.New("user has been band")
	ErrUnsupportedMsgType = errors.New("unsupported msg type")
)

type EventReply struct {
	Reply    string `json:"reply"`
	ATSender bool   `json:"at_sender"`
}

const (
	MTPrivate = "private"
	MTGroup   = "group"

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

func (uc *EventUseCase) Event(ctx context.Context, m *Message) (reply EventReply, err error) {
	if err := ratelimit.Rate(ctx, "cq_event", strconv.Itoa(int(m.UserID))); err != nil {
		return reply, err
	}

	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		if err := uc.repo.SaveUser(ctx, tx, m.User); err != nil {
			return err
		}
		return uc.repo.SaveMessage(ctx, tx, m)
	}); err != nil {
		logrus.Errorf("event message save error %v\n", err)
	}

	// User Filter, 检查是否被禁用
	if m.User.Ban {
		return reply, errors.Wrapf(ErrUserBand, "user: %#v", m.User)
	}

	// 排除类型
	if m.MessageType != MTPrivate && m.MessageType != MTGroup {
		return reply, errors.Wrapf(ErrUnsupportedMsgType, "type: %s", m.MessageType)
	}

	// 排除群号
	if m.MessageType == MTGroup {
		if e := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
			_, e := uc.repo.GetGroupByID(ctx, tx, m.GroupID.Int64)
			return e
		}); e != nil {
			// 只显示主要原因, 抛弃栈信息
			return reply, errors.WithMessage(errors.Cause(e), strconv.Itoa(int(m.GroupID.Int64)))
		}
		reply.ATSender = true
	}

	reply.Reply, err = uc.doEvent(ctx, m)
	if err != nil {
		return reply, err
	}

	err = database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		return uc.repo.SetMessageReply(ctx, tx, m.ID, reply.Reply)
	})

	return reply, err
}

func (uc *EventUseCase) doEvent(ctx context.Context, m *Message) (string, error) {
	// 等待3秒 回复
	defer time.Sleep(time.Second * 2)
	defer logrus.Infof("do event reply %s\n", m.Reply)

	// 匹配简单回复
	return autoreply.Reply(ctx, &autoreply.Message{
		UserID:  strconv.Itoa(int(m.UserID)),
		Message: m.Message,
	})
}

func Equal(s, v string) bool {
	return strings.ReplaceAll(s, " ", "") == strings.ReplaceAll(v, " ", "")
}

func Prefix(s, v string) bool {
	return strings.HasPrefix(strings.ReplaceAll(s, " ", ""), v)
}
