package v1

import "github.com/moyrne/tebot/internal/models"

type QMessage struct {
	ID          int    `json:"id"`
	Time        int    `json:"time"`         // 时间戳
	SelfID      int    `json:"self_id"`      // 用户ID
	PostType    string `json:"post_type"`    // 上报类型	message
	MessageType string `json:"message_type"` // 消息类型	private
	SubType     string `json:"sub_type"`     // 消息子类型 friend,group,group_self,other
	TempSource  int    `json:"temp_source"`  // 临时会话来源
	MessageID   int    `json:"message_id"`   // 消息 ID
	UserID      int    `json:"user_id"`      // 发送者 QQ 号
	Message     string `json:"message"`      // 消息内容
	RawMessage  string `json:"raw_message"`  // 原始消息内容
	Font        int    `json:"font"`         // 字体
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

func (s TempSource) String() string {
	switch s {
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
		return "UnKnow"
	}
}
