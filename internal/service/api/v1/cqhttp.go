package v1

import "github.com/gin-gonic/gin"

type CqHTTP struct{}

const (
	STFriend    = "friend"
	STGroup     = "group"
	STGroupSelf = "group_self"
	STOther     = "other"
)

type QMessage struct {
	ID          int    `json:"id"` // 时间戳
	Time        int    `json:"time"`
	SelfID      int    `json:"self_id"`
	PostType    string `json:"post_type"`    // 上报类型	message
	MessageType string `json:"message_type"` // 消息类型	private
	SubType     string `json:"sub_type"`     // 消息子类型 friend,group,group_self,other
	TempSource  int    `json:"temp_source"`  // 临时会话来源
	/*
		0	群聊
		1	QQ咨询
		2	查找
		3	QQ电影
		4	热聊
		6	验证消息
		7	多人聊天
		8	约会
		9	通讯录
	*/
	MessageID  int    `json:"message_id"`  // 消息ID
	UserID     int    `json:"user_id"`     // 发送者 QQ 号
	Message    string `json:"message"`     // 消息内容
	RawMessage string `json:"raw_message"` // 原始消息内容
	Font       int    `json:"font"`        // 字体
}

// HTTP 文档 https://github.com/ishkong/go-cqhttp-docs/tree/main/docs/event
func (h CqHTTP) HTTP(c *gin.Context) {
	var params QMessage
	if err := c.BindJSON(&params); err != nil {
		// TODO log error
		return
	}

	switch params.SubType {
	case STFriend:
		h.friend(c, params)
	case STGroup, STGroupSelf:
		h.group(c, params)
	case STOther:
		fallthrough
	default:
		// TODO log error
		return
	}
}

func (h CqHTTP) friend(c *gin.Context, params QMessage) {

}

func (h CqHTTP) group(c *gin.Context, params QMessage) {

}
