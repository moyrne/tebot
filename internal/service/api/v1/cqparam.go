package v1

import (
	"database/sql"
	"github.com/moyrne/tebot/internal/models"
)

type QMessage struct {
	ID          int         `json:"id"`
	Time        int         `json:"time"`         // 时间戳
	SelfID      int         `json:"self_id"`      // 用户ID
	PostType    string      `json:"post_type"`    // 上报类型	meta_event, message
	MessageType string      `json:"message_type"` // 消息类型	private, group
	SubType     string      `json:"sub_type"`     // 消息子类型 [private]private, group, group_self, other; [private] normal, anonymous, notice
	TempSource  *TempSource `json:"temp_source"`  // 临时会话来源
	MessageID   int         `json:"message_id"`   // 消息 ID
	GroupID     int64       `json:"group_id"`     // 群ID
	UserID      int         `json:"user_id"`      // 发送者 QQ 号
	Message     string      `json:"message"`      // 消息内容
	RawMessage  string      `json:"raw_message"`  // 原始消息内容
	Font        int         `json:"font"`         // 字体
	Sender      QSender     `json:"sender"`       // 发送人信息
}

func (m QMessage) Model() *models.QMessage {
	return &models.QMessage{
		ID:          int64(m.ID),
		Time:        m.Time,
		SelfID:      m.SelfID,
		PostType:    m.PostType,
		MessageType: m.MessageType,
		SubType:     m.SubType,
		TempSource:  m.TempSource.String(),
		MessageID:   m.MessageID,
		GroupID: sql.NullInt64{
			Int64: m.GroupID,
		},
		UserID:     m.UserID,
		Message:    m.Message,
		RawMessage: m.RawMessage,
		Font:       m.Font,
		QUser: &models.QUser{
			QUID:     m.Sender.UserID,
			Nickname: m.Sender.Nickname,
			Sex:      m.Sender.Sex,
			Age:      m.Sender.Age,
		},
	}
}

type QSender struct {
	UserID   int    `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 昵称
	Sex      string `json:"sex"`      // 性别, male 或 female 或 unknown
	Age      int    `json:"age"`      // 年龄
}

type TempSource int

const (
	TSGroup       = 0 // 群聊
	TSConsult     = 1 // QQ咨询
	TSSearch      = 2 // 查找
	TSFilm        = 3 // QQ电影
	TSHotTalk     = 4 // 热聊
	TSVerify      = 6 // 验证消息
	TSMultiChat   = 7 // 多人聊天
	TSAppointment = 8 // 约会
	TSMailList    = 9 // 通讯录
)

func (s *TempSource) String() string {
	if s == nil {
		return ""
	}
	switch *s {
	case TSGroup:
		return models.TSGroup // 群聊
	case TSConsult:
		return models.TSConsult // QQ咨询
	case TSSearch:
		return models.TSSearch // 查找
	case TSFilm:
		return models.TSFilm // QQ电影
	case TSHotTalk:
		return models.TSHotTalk // 热聊
	case TSVerify:
		return models.TSVerify // 验证消息
	case TSMultiChat:
		return models.TSMultiChat // 多人聊天
	case TSAppointment:
		return models.TSAppointment // 约会
	case TSMailList:
		return models.TSMailList // 通讯录
	default:
		return "unknown"
	}
}

type Reply struct {
	Reply       string `json:"reply"`                  // 要回复的内容
	AutoEscape  bool   `json:"auto_escape,omitempty"`  // 消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 reply 字段是字符串时有效
	ATSender    bool   `json:"at_sender,omitempty"`    // 是否要在回复开头 at 发送者 ( 自动添加 ) , 发送者是匿名用户时无效	at 发送者
	Delete      bool   `json:"delete,omitempty"`       // 撤回该条消息	不撤回
	Kick        bool   `json:"kick,omitempty"`         // 把发送者踢出群组 ( 需要登录号权限足够 ) , 不拒绝此人后续加群请求, 发送者是匿名用户时无效	不踢
	Ban         bool   `json:"ban,omitempty"`          // 把发送者禁言 ban_duration 指定时长, 对匿名用户也有效	不禁言
	BanDuration int    `json:"ban_duration,omitempty"` // 禁言时长	30 分钟
}
