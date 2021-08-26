package biz

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/pkg/autoreply"
	"github.com/moyrne/tebot/internal/pkg/logs"
	"github.com/moyrne/tractor/dbx"
)

// SyncReply 同步 回复
func SyncReply(ctx context.Context, repo ReplyRepo) {
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			delaySync(ctx, repo)
		}
	}()
}

func delaySync(ctx context.Context, repo ReplyRepo) {
	// 5分钟同步一次
	defer time.Sleep(time.Minute * 5)
	var replies []Reply
	if err := database.NewTransaction(ctx, func(ctx context.Context, tx dbx.Transaction) (err error) {
		replies, err = repo.Replies(ctx, tx)
		return err
	}); err != nil {
		logs.Error("sync reply", "error", err)
		return
	}

	repliesMap := autoreply.Replies{}
	for _, reply := range replies {
		userID := strconv.Itoa(int(reply.UserID))
		if reply.UserID == 0 {
			userID = autoreply.DefaultReply
		}
		if _, ok := repliesMap[userID]; !ok {
			repliesMap[userID] = map[string]autoreply.ReplyRow{}
		}
		var r []string
		if err := json.Unmarshal([]byte(reply.Replies), &r); err != nil {
			logs.Error("delay sync unmarshal", "data", reply, "error", err)
			continue
		}
		repliesMap[userID][reply.Msg] = autoreply.ReplyRow{
			Msg:     reply.Msg,
			Replies: r,
		}
	}
	// copy on write
	autoreply.RefreshReplies(repliesMap)
}
