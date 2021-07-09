package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/moyrne/tebot/internal/analyze"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/models"
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
		// 忽略 心跳检测
		return
	}

	if params.PostType != PTMessage {
		// 忽略 非消息事件
		logs.Error("unknown params", "post_type", params.PostType)
		return
	}

	qmModel := params.Model()
	// 优先提供服务 记录失败忽略
	err := database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx *sqlx.Tx) error {
		return qmModel.Insert(c.Request.Context(), tx)
	})
	if err != nil {
		logs.Error("insert CqHTTP params failed", "error", err)
	}

	// User Filter, 检查是否被禁用
	if qmModel.QUser.Ban {
		logs.Info("quser has been band", "user", qmModel.QUser)
		return
	}
	var reply Reply
	switch params.MessageType {
	case MTPrivate:
		reply, err = h.private(c, qmModel)
		if err != nil {
			logs.Info("private", "error", err)
			return
		}
	case MTGroup:
		reply, err = h.group(c, qmModel)
		if err != nil {
			logs.Info("group", "error", err)
			return
		}
	default:
		// log error
		logs.Error("unsupported message_type", "type", params.MessageType)
		return
	}

	if err = database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx *sqlx.Tx) error {
		qmModel.Reply = reply.Reply
		return qmModel.SetReply(ctx, tx)
	}); err != nil {
		logs.Error("update cqhttp params failed", "error", err)
	}
	c.JSON(http.StatusOK, reply)
}

func (h CqHTTP) private(c *gin.Context, params *models.QMessage) (Reply, error) {
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return Reply{}, err
	}
	return Reply{Reply: reply}, nil
}

func (h CqHTTP) group(c *gin.Context, params *models.QMessage) (Reply, error) {
	if err := database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := models.GetQGroupByQGID(ctx, tx, params.GroupID)
		return err
	}); err != nil {
		return Reply{}, err
	}
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return Reply{}, err
	}
	return Reply{Reply: reply, ATSender: true}, nil
}
