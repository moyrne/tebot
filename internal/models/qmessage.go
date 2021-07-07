package models

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

// QMessage QQ 聊天记录
type QMessage struct {
	ID          int    `json:"id"` // 时间戳
	Time        int    `json:"time"`
	SelfID      int    `json:"self_id"`
	PostType    string `json:"post_type"`    // 上报类型	message
	MessageType string `json:"message_type"` // 消息类型	private
	SubType     string `json:"sub_type"`     // 消息子类型 friend,group,group_self,other
	TempSource  string `json:"temp_source"`  // 临时会话来源
	/*
		0->group		群聊
		1->consult		QQ咨询
		2->search		查找
		3->film			QQ电影
		4->hottalk		热聊
		6->verify		验证消息
		7->multichat	多人聊天
		8->appointment	约会
		9->maillist		通讯录
	*/
	MessageID  int    `json:"message_id"`  // 消息ID
	UserID     int    `json:"user_id"`     // 发送者 QQ 号
	Message    string `json:"message"`     // 消息内容
	RawMessage string `json:"raw_message"` // 原始消息内容
	Font       int    `json:"font"`        // 字体
}

func (m QMessage) TableName() string {
	return "q_message"
}
