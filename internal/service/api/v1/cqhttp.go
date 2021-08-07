package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/internal/analyze"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/models"
	"github.com/moyrne/tebot/internal/service/commands"
	"github.com/moyrne/tractor/dbx"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
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
		logs.Info("cqhttp", "unmarshal error", err)
		return
	}

	if params.PostType == PTEvent {
		// 忽略 心跳检测
		commands.CQHeartBeat()
		return
	}

	if params.PostType != PTMessage {
		// 忽略 非消息事件
		logs.Error("unknown params", "params", params)
		return
	}

	qmModel := params.Model()
	// 优先提供服务 记录失败忽略
	err := database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx dbx.Transaction) error {
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
	case MTGroup:
		reply, err = h.group(c, qmModel)
	default:
		// log error
		logs.Error("unsupported message_type", "type", params.MessageType)
		return
	}

	if errors.Is(err, analyze.ErrNotMatch) {
		return
	}
	if err != nil {
		logs.Info("get reply", "error", err)
		return
	}

	if err = database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx dbx.Transaction) error {
		qmModel.Reply = reply.Reply
		return qmModel.SetReply(ctx, tx)
	}); err != nil {
		logs.Error("update cqhttp params failed", "error", err)
	}
	// 等待3秒 回复
	time.Sleep(time.Second * 3)
	logs.Info("reply", "content", reply)
	c.JSON(http.StatusOK, reply)
}

func (h CqHTTP) private(c *gin.Context, params *models.QMessage) (Reply, error) {
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err == nil {
		return Reply{Reply: reply}, nil
	}
	return Reply{}, err
}

func (h CqHTTP) group(c *gin.Context, params *models.QMessage) (Reply, error) {
	if err := database.NewTransaction(c.Request.Context(), func(ctx context.Context, tx dbx.Transaction) error {
		_, err := models.GetQGroupByQGID(ctx, tx, params.GroupID)
		return err
	}); err != nil {
		// 只显示主要原因, 抛弃栈信息
		return Reply{}, errors.WithMessage(errors.Cause(err), strconv.Itoa(int(params.GroupID.Int64)))
	}
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		return Reply{}, err
	}
	return Reply{Reply: reply, ATSender: true}, nil
}
