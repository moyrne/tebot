package cqhttp

import (
	"context"
	"database/sql"

	api "github.com/moyrne/tebot/api/cqhttp"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
)

// TODO 定义标准 const
const (
	PTMessage = "message"
	PTEvent   = "meta_event"
	MTPrivate = "private"
	MTGroup   = "group"
)

type Server struct {
	biz *cqhttp.EventUseCase
}

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

	// 优先提供服务 记录失败忽略
	reply, err := s.biz.Event(ctx, ToQMessage(m))
	return api.Reply{Reply: reply.Reply, ATSender: reply.ATSender}, err
}

func ToQMessage(m api.QMessage) *cqhttp.Message {
	return &cqhttp.Message{
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
		User: &cqhttp.User{
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
		return cqhttp.MTPrivate
	case MTGroup:
		return cqhttp.MTGroup
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
		return cqhttp.TSGroup // 群聊
	case api.TSConsult:
		return cqhttp.TSConsult // QQ咨询
	case api.TSSearch:
		return cqhttp.TSSearch // 查找
	case api.TSFilm:
		return cqhttp.TSFilm // QQ电影
	case api.TSHotTalk:
		return cqhttp.TSHotTalk // 热聊
	case api.TSVerify:
		return cqhttp.TSVerify // 验证消息
	case api.TSMultiChat:
		return cqhttp.TSMultiChat // 多人聊天
	case api.TSAppointment:
		return cqhttp.TSAppointment // 约会
	case api.TSMailList:
		return cqhttp.TSMailList // 通讯录
	default:
		return "unknown"
	}
}
