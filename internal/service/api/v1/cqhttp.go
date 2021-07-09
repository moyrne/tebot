package v1

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/internal/analyze"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/pkg/errors"
	"net/http"
)

type CqHTTP struct{}

const (
	PTMessage = "message"
	PTEvent   = "meta_event"
	MTPrivate = "private"
	MTGroup   = "group"
)

// HTTP 文档 https://github.com/ishkong/go-cqhttp-docs/tree/main/docs/event
func (h CqHTTP) HTTP(c *gin.Context) {
	var params QMessage
	if err := c.BindJSON(&params); err != nil {
		// TODO log error
		return
	}

	// TODO 限流 防止封号 (20sCD)

	if params.PostType == PTEvent {
		// 忽略
		return
	}

	if params.PostType != PTMessage {
		// 忽略
		logs.Error("unknown params", "post_type", params.PostType)
		return
	}

	// 优先提供服务 记录失败忽略
	err := database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx *sql.Tx) error {
		return params.Model().Insert(c.Request.Context(), tx)
	})
	if err != nil {
		logs.Error("insert cqhttp params failed", "error", err)
	}

	var reply Reply
	switch params.MessageType {
	case MTPrivate:
		reply, err = h.private(c, params)
		if err != nil {
			logs.Info("private", "error", err)
			return
		}
	case MTGroup:
		reply, err = h.group(c, params)
		if err != nil {
			logs.Info("group", "error", err)
			return
		}
	default:
		// TODO log error
		logs.Error("unsupported message_type", "type", params.MessageType)
		return
	}
	if err = database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx *sql.Tx) error {
		return params.Model().SetReply(ctx, tx)
	}); err != nil {
		logs.Error("insert cqhttp params failed", "error", err)
	}
	c.JSON(http.StatusOK, reply)
	return
}

func (h CqHTTP) private(c *gin.Context, params QMessage) (Reply, error) {
	// TODO User Filter, 检查是否有权限
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return Reply{}, errors.WithStack(err)
	}
	return Reply{Reply: reply}, nil
}

func (h CqHTTP) group(c *gin.Context, params QMessage) (Reply, error) {
	// TODO Group Filter + User Filter, 检查是否有权限
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return Reply{}, errors.WithStack(err)
	}
	return Reply{Reply: reply, ATSender: true}, nil
}
