package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/internal/analyze"
	"github.com/moyrne/tebot/internal/logs"
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

	switch params.MessageType {
	case MTPrivate:
		h.private(c, params)
	case MTGroup:
		h.group(c, params)
	default:
		// TODO log error
		return
	}
}

func (h CqHTTP) private(c *gin.Context, params QMessage) {
	// TODO User Filter, 检查是否有权限
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		logs.Info("private analyze.Analyze", "error", err)
		return
	}
	c.JSON(http.StatusOK, Reply{Reply: reply})
}

func (h CqHTTP) group(c *gin.Context, params QMessage) {
	// TODO Group Filter + User Filter, 检查是否有权限
	reply, err := analyze.Analyze(c.Request.Context(), analyze.Params{QUID: params.UserID, Message: params.Message})
	if err != nil {
		logs.Info("group analyze.Analyze", "error", err)
		return
	}
	c.JSON(http.StatusOK, Reply{Reply: reply, ATSender: true})
}
