package autoreply

import (
	"context"
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/pkg/errors"
)

const DefaultReply = "default"

type Message struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type (
	Match    func(value, match string) bool
	Function func(ctx context.Context, m *Message) (string, error)
)

var (
	matches   = map[string]Match{}
	functions = map[string]Function{}
)

func RegisterMatches(name string, match Match) {
	matches[name] = match
}

func RegisterFunctions(name string, function Function) {
	functions[name] = function
}

type ReplyRow struct {
	Msg     string `json:"msg"`
	Matches string `json:"matches"`

	Function string   `json:"function"`
	Replies  []string `json:"repliesValue"`
}

type Replies map[string]map[string]ReplyRow

var repliesValue atomic.Value

func init() {
	repliesValue.Store(Replies{})
}

func RefreshReplies(r Replies) {
	repliesValue.Store(r)
}

func loadReplies() Replies {
	return repliesValue.Load().(Replies)
}

var ErrNotMatch = errors.New("not match")

// Reply 简单回复
func Reply(ctx context.Context, m *Message) (string, error) {
	replies := loadReplies()
	custom, ok := replies[m.UserID]
	if ok {
		resp, err := rangeReply(ctx, custom, m)
		if err == nil {
			return resp, nil
		}
		if !errors.Is(err, ErrNotMatch) {
			return "", err
		}
	}

	return rangeReply(ctx, replies[DefaultReply], m)
}

func rangeReply(ctx context.Context, replies map[string]ReplyRow, m *Message) (string, error) {
	for msg, reply := range replies {
		mts := strings.Split(reply.Matches, ",")
		for _, match := range mts {
			if matches[match](m.Message, msg) {
				if fn, ok := functions[reply.Function]; ok {
					return fn(ctx, m)
				}
				return randReply(reply.Replies), nil
			}
		}
	}
	return "", errors.WithStack(ErrNotMatch)
}

func randReply(replies []string) string {
	rd := rand.Intn(len(replies))
	return replies[rd]
}
