package cqhttp

import (
	"context"
	"database/sql"
	api "github.com/moyrne/tebot/api/cqhttp"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
	"strconv"
	"time"

	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

const (
	PTMessage = "message"
	PTEvent   = "meta_event"
	MTPrivate = "private"
	MTGroup   = "group"
)

type Server struct{}

// Event 文档 https://github.com/ishkong/go-cqhttp-docs/tree/main/docs/event
func (s Server) Event(ctx context.Context, m api.QMessage) (api.Reply, error) {
	if m.PostType == PTEvent {
		// 忽略 心跳检测
		keepalive.CQHeartBeat()
		return api.Reply{}, nil
	}

	if m.PostType != PTMessage {
		// 忽略 非消息事件
		logs.Error("unknown params", "params", m)
		return api.Reply{}, nil
	}

	qmModel := ToQMessage(m)
	// 优先提供服务 记录失败忽略
	err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		return qmModel.Insert(ctx, tx)
	})
	if err != nil {
		logs.Error("insert CqHTTP params failed", "error", err)
	}

	// User Filter, 检查是否被禁用
	if qmModel.QUser.Ban {
		logs.Info("quser has been band", "user", qmModel.QUser)
		return api.Reply{}, nil
	}
	var reply api.Reply
	switch m.MessageType {
	case MTPrivate:
		reply, err = s.private(ctx, qmModel)
	case MTGroup:
		reply, err = s.group(ctx, qmModel)
	default:
		// log error
		logs.Error("unsupported message_type", "type", m.MessageType)
		return api.Reply{}, nil
	}

	if errors.Is(err, cqhttp.ErrNotMatch) {
		return api.Reply{}, nil
	}
	if err != nil {
		logs.Info("get reply", "error", err)
		return api.Reply{}, nil
	}

	if err = database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		qmModel.Reply = reply.Reply
		return qmModel.SetReply(ctx, tx)
	}); err != nil {
		logs.Error("update cqhttp params failed", "error", err)
	}
	// 等待3秒 回复
	time.Sleep(time.Second * 3)
	logs.Info("reply", "content", reply)
	return reply, nil
}
func (s Server) private(ctx context.Context, params *data.QMessage) (api.Reply, error) {
	reply, err := cqhttp.Analyze(ctx, cqhttp.Params{QUID: params.UserID, Message: params.Message})
	if err == nil {
		return api.Reply{Reply: reply}, nil
	}
	return api.Reply{}, err
}

func (s Server) group(ctx context.Context, params *data.QMessage) (api.Reply, error) {
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) error {
		_, err := data.GetQGroupByQGID(ctx, tx, params.GroupID)
		return err
	}); err != nil {
		// 只显示主要原因, 抛弃栈信息
		return api.Reply{}, errors.WithMessage(errors.Cause(err), strconv.Itoa(int(params.GroupID.Int64)))
	}
	reply, err := cqhttp.Analyze(ctx, cqhttp.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return api.Reply{}, err
	}
	return api.Reply{Reply: reply, ATSender: true}, nil
}

func ToQMessage(m api.QMessage) *data.QMessage {
	return &data.QMessage{
		ID:          int64(m.ID),
		Time:        m.Time,
		SelfID:      m.SelfID,
		PostType:    m.PostType,
		MessageType: m.MessageType,
		SubType:     m.SubType,
		TempSource:  ToTempSource(m.TempSource),
		MessageID:   m.MessageID,
		GroupID: sql.NullInt64{
			Int64: m.GroupID,
		},
		UserID:     m.UserID,
		Message:    m.Message,
		RawMessage: m.RawMessage,
		Font:       m.Font,
		QUser: &data.QUser{
			QUID:     m.Sender.UserID,
			Nickname: m.Sender.Nickname,
			Sex:      m.Sender.Sex,
			Age:      m.Sender.Age,
		},
	}
}

func ToTempSource(s *api.TempSource) string {
	if s == nil {
		return ""
	}
	switch *s {
	case api.TSGroup:
		return data.TSGroup // 群聊
	case api.TSConsult:
		return data.TSConsult // QQ咨询
	case api.TSSearch:
		return data.TSSearch // 查找
	case api.TSFilm:
		return data.TSFilm // QQ电影
	case api.TSHotTalk:
		return data.TSHotTalk // 热聊
	case api.TSVerify:
		return data.TSVerify // 验证消息
	case api.TSMultiChat:
		return data.TSMultiChat // 多人聊天
	case api.TSAppointment:
		return data.TSAppointment // 约会
	case api.TSMailList:
		return data.TSMailList // 通讯录
	default:
		return "unknown"
	}
}
