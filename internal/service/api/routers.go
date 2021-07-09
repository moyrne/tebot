package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/moyrne/tebot/internal/service/api/v1"
)

func NewRouter() *gin.Engine {
	e := gin.Default()
	var h v1.CqHTTP
	e.POST("/", h.HTTP)
	return e
}
