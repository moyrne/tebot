package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
	v1 "github.com/moyrne/tebot/internal/service/cqhttp"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/moyrne/tebot/configs"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
)

func main() {
	if err := configs.LoadConfig(); err != nil {
		log.Fatalln("read config error", err)
	}
	writer, err := logs.FileWriter()
	if err != nil {
		log.Fatalln("new file writer error", err)
	}
	defer writer.Close()
	logs.Init(writer)
	if err := database.ConnectMySQL(); err != nil {
		logs.Panic("db connect", "error", err)
	}
	cqhttp.SyncReply(context.Background())
	if err := database.ConnectRedis(); err != nil {
		logs.Panic("redis connect", "error", err)
	}
	go keepalive.StartCQHTTP()
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	var h v1.CqHTTP
	cqhttp.InitLimiter()
	e.POST("/", h.HTTP)
	if err := e.Run("127.0.0.1:7771"); err != nil {
		logs.Panic("service run", "error", err)
	}
}
