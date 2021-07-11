package analyze

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/models"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func PrintMenu(_ context.Context, _ Params) (string, error) {
	return "๑ 菜单\n" +
		"๑ 1.绑定位置;(绑定位置 深圳)\n" +
		"๑ 2.签到", nil
}

type QReplyRow struct {
	Msg    string   `json:"msg"`
	Weight int      `json:"weight"`
	Reply  []string `json:"reply"`
}

var (
	simpleMu    sync.RWMutex
	simpleReply map[int]map[string]QReplyRow
)

// SimpleReply 简单回复
func SimpleReply(_ context.Context, params Params) (string, error) {
	simpleMu.RLock()
	defer simpleMu.RUnlock()
	custom, ok := simpleReply[params.QUID]
	if ok {
		for msg, reply := range custom {
			if strings.Contains(params.Message, msg) {
				rd := rand.Intn(len(reply.Reply) - 1)
				return reply.Reply[rd+1], nil
			}
		}
	}

	def := simpleReply[0]
	for msg, reply := range def {
		if strings.Contains(params.Message, msg) {
			rd := rand.Intn(len(reply.Reply) - 1)
			return reply.Reply[rd+1], nil
		}
	}

	return "", errors.WithStack(ErrNotMatch)
}

// SyncReply 同步 回复
func SyncReply(ctx context.Context) {
	go func() {
		for {
			delaySync(ctx)
		}
	}()
}

func delaySync(ctx context.Context) {
	// 5分钟同步一次
	defer time.Sleep(time.Minute * 5)
	var replies []models.QReply
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) (err error) {
		replies, err = models.SelectQReply(ctx, tx)
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
		simpleReply[reply.QUID][reply.Msg] = QReplyRow{
			Msg:    reply.Msg,
			Weight: reply.Weight,
			Reply:  reply.Reply,
		}
	}
}
