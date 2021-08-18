package cqhttp

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
)

type ReplyRow struct {
	Msg     string `json:"msg"`
	Matches string `json:"matches"`

	Function string   `json:"function"`
	Replies  []string `json:"replies"`
}

var (
	simpleMu    sync.RWMutex
	simpleReply map[int64]map[string]ReplyRow
)

// ReplyMethod 简单回复
func ReplyMethod(ctx context.Context, uc *EventUseCase, m *Message) (string, error) {
	simpleMu.RLock()
	defer simpleMu.RUnlock()
	custom, ok := simpleReply[m.UserID]
	if ok {
		resp, err := rangeSimpler(ctx, uc, custom, m)
		if err == nil {
			return resp, nil
		}
		if !errors.Is(err, ErrNotMatch) {
			return "", err
		}
	}

	return rangeSimpler(ctx, uc, simpleReply[0], m)
}

func rangeSimpler(ctx context.Context, uc *EventUseCase, replies map[string]ReplyRow, m *Message) (string, error) {
	for msg, reply := range replies {
		matches := strings.Split(reply.Matches, ",")
		for _, match := range matches {
			if Matches[match].Fn(m.Message, msg) {
				if err := rateLimiter.Rate(ctx, "reply", m.UserID); err != nil {
					return "", err
				}
				if fn, ok := Functions[reply.Function]; ok {
					return fn.Fn(ctx, uc, m)
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

// SyncReply 同步 回复
func SyncReply(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			delaySync(ctx)
		}
	}()
}

func delaySync(ctx context.Context) {
	// 5分钟同步一次
	defer time.Sleep(time.Minute * 5)
	var replies []data.Reply
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) (err error) {
		replies, err = data.SelectQReply(ctx, tx)
		return err
	}); err != nil {
		logs.Error("sync reply", "error", err)
		return
	}

	simpleMu.Lock()
	defer simpleMu.Unlock()
	simpleReply = map[int64]map[string]ReplyRow{}
	for _, reply := range replies {
		if _, ok := simpleReply[reply.UserID]; !ok {
			simpleReply[reply.UserID] = map[string]ReplyRow{}
		}
		var r []string
		if err := json.Unmarshal([]byte(reply.Reply), &r); err != nil {
			logs.Error("delay sync unmarshal", "data", reply, "error", err)
			continue
		}
		simpleReply[reply.UserID][reply.Msg] = ReplyRow{
			Msg:     reply.Msg,
			Replies: r,
		}
	}
}
