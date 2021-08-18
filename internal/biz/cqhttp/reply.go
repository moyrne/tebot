package cqhttp

type Reply struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Msg     string `json:"msg"`     // 前缀搜索
	Matches string `json:"matches"` // TODO 可保存多个，使用','隔开

	Function string `json:"function"`

	Replies string `json:"replies"`
}
