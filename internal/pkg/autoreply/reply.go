package autoreply

import "context"

type Message struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type Function func(ctx context.Context, m *Message) (string, error)

var (
	matches   = map[string]Function{}
	functions = map[string]Function{}
)

func RegisterMatches(name string, match Function) {
	matches[name] = match
}

func RegisterFunctions(name string, function Function) {
	functions[name] = function
}
