package service

import (
	"context"
	"database/sql"

	"github.com/moyrne/tebot/api"
	"github.com/moyrne/tebot/internal/biz"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
	"github.com/moyrne/tebot/internal/pkg/logs"
	"github.com/sirupsen/logrus"
)

// TODO 定义标准 const
const (
	PTMessage = "message"
	PTEvent   = "meta_event"
	MTPrivate = "private"
	MTGroup   = "group"
)

func NewEventServer(repo biz.EventRepo) EventServer {
	return EventServer{biz: biz.NewEventUseCase(repo)}
}

type EventServer struct {
	biz *biz.EventUseCase
}

// Event 文档 https://github.com/ishkong/go-cqhttp-docs/tree/main/docs/event
func (s EventServer) Event(ctx context.Context, m api.QMessage) (api.Reply, error) {
	if m.PostType == PTEvent {
		// 忽略 心跳检测
		keepalive.CQHeartBeat()
		return api.Reply{}, nil
	}

	if m.PostType != PTMessage {
		// 忽略 非消息事件
		logrus.Errorf("unknown params %s\n", logs.JSONMarshalIgnoreErr(m))
		return api.Reply{}, nil
	}

	// 优先提供服务 记录失败忽略
	reply, err := s.biz.Event(ctx, ToQMessage(m))
	return api.Reply{Reply: reply.Reply, ATSender: reply.ATSender}, err
}

func ToQMessage(m api.QMessage) *biz.Message {
	return &biz.Message{
		ID:          int64(m.ID),
		Time:        m.Time,
		SelfID:      m.SelfID,
		PostType:    m.PostType,
		MessageType: ToMessageType(m.MessageType),
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
		User: &biz.User{
			UserID:   m.Sender.UserID,
			Nickname: m.Sender.Nickname,
			Sex:      m.Sender.Sex,
			Age:      m.Sender.Age,
		},
	}
}

func ToMessageType(t string) string {
	switch t {
	case MTPrivate:
		return biz.MTPrivate
	case MTGroup:
		return biz.MTGroup
	default:
		return "unknown"
	}
}

func ToTempSource(s *api.TempSource) string {
	if s == nil {
		return ""
	}
	switch *s {
	case api.TSGroup:
		return biz.TSGroup // 群聊
	case api.TSConsult:
		return biz.TSConsult // QQ咨询
	case api.TSSearch:
		return biz.TSSearch // 查找
	case api.TSFilm:
		return biz.TSFilm // QQ电影
	case api.TSHotTalk:
		return biz.TSHotTalk // 热聊
	case api.TSVerify:
		return biz.TSVerify // 验证消息
	case api.TSMultiChat:
		return biz.TSMultiChat // 多人聊天
	case api.TSAppointment:
		return biz.TSAppointment // 约会
	case api.TSMailList:
		return biz.TSMailList // 通讯录
	default:
		return "unknown"
	}
}
