package api

import (
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/internal/analyze"
	v1 "github.com/moyrne/tebot/internal/service/api/v1"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	var h v1.CqHTTP
	analyze.InitLimiter()
	e.POST("/", h.HTTP)
	return e
}
