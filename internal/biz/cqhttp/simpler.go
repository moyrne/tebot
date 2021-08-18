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

func PrintMenu(_ context.Context, _ *QMessage) (string, error) {
	return "๑ 菜单\n" +
		"๑ 1.绑定位置;(绑定位置 深圳)\n" +
		"๑ 2.签到", nil
}

type QReplyRow struct {
	Msg   string   `json:"msg"`
	Reply []string `json:"reply"`
}

var (
	simpleMu    sync.RWMutex
	simpleReply map[int]map[string]QReplyRow
)

// SimpleReply 简单回复
func SimpleReply(ctx context.Context, m *QMessage) (string, error) {
	simpleMu.RLock()
	defer simpleMu.RUnlock()
	custom, ok := simpleReply[m.UserID]
	if ok {
		resp, err := rangeSimpler(ctx, custom, m, func(s, v string) bool {
			return strings.Contains(s, v)
		})
		if err == nil {
			return resp, nil
		}
		if !errors.Is(err, ErrNotMatch) {
			return "", err
		}
	}

	return rangeSimpler(ctx, simpleReply[0], m, func(s, v string) bool {
		return strings.Contains(s, v)
	})
}

func rangeSimpler(ctx context.Context, replies map[string]QReplyRow, m *QMessage, match func(s, v string) bool) (string, error) {
	for msg, reply := range replies {
		if match(m.Message, msg) {
			if err := rateLimiter.Rate(ctx, "simple", m.UserID); err != nil {
				return "", err
			}
			return randReply(reply.Reply), nil
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
	var replies []data.QReply
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) (err error) {
		replies, err = data.SelectQReply(ctx, tx)
		return err
	}); err != nil {
		logs.Error("sync reply", "error", err)
		return
	}

	simpleMu.Lock()
	defer simpleMu.Unlock()
	simpleReply = map[int]map[string]QReplyRow{}
	for _, reply := range replies {
		if _, ok := simpleReply[reply.QUID]; !ok {
			simpleReply[reply.QUID] = map[string]QReplyRow{}
		}
		var r []string
		if err := json.Unmarshal([]byte(reply.Reply), &r); err != nil {
			logs.Error("delay sync unmarshal", "data", reply, "error", err)
			continue
		}
		simpleReply[reply.QUID][reply.Msg] = QReplyRow{
			Msg:   reply.Msg,
			Reply: r,
		}
	}
}
