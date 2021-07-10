package analyze

import "context"

func PrintMenu(_ context.Context, _ Params) (string, error) {
	return "[菜单]\n" +
		"1.绑定位置;(绑定位置 深圳)\n" +
		"2.签到", nil
}

// SimpleReply 简单回复
func SimpleReply(ctx context.Context, params Params) (string, error) {
	return "", nil
}
