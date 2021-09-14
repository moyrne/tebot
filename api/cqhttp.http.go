package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CQServer interface {
	Event(ctx context.Context, message QMessage) (Reply, error)
}

func RegisterCQServer(e *gin.Engine, server CQServer) {
	e.POST("/", func(c *gin.Context) {
		var message QMessage
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "json unmarshal failed"})
			return
		}
		reply, err := server.Event(c.Request.Context(), message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, reply)
	})
}
