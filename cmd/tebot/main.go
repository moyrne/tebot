package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/api"
	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
	"github.com/moyrne/tebot/internal/service"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/moyrne/tebot/configs"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
)

func main() {
	// 加载配置文件
	if err := configs.LoadConfig(); err != nil {
		log.Fatalln("read config error", err)
	}

	// 初始化日志文件
	writer, err := logs.FileWriter()
	if err != nil {
		log.Fatalln("new file writer error", err)
	}
	defer writer.Close()

	// 连接数据库
	if err := database.ConnectMySQL(); err != nil {
		logs.Panic("db connect", "error", err)
	}

	// 连接Redis
	if err := database.ConnectRedis(); err != nil {
		logs.Panic("redis connect", "error", err)
	}

	// 初始化日志
	logs.Init(writer, data.NewLogRepo())

	// 心跳检测
	go keepalive.StartCQHTTP()

	// 初始化限流
	cqhttp.InitLimiter()

	// 启动自动回复 同步
	cqhttp.SyncReply(context.Background(), data.NewReplyRepo())

	// 启动WEB服务
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	api.RegisterServer(e, service.NewEventServer(data.NewEventRepo()))
	if err := e.Run("127.0.0.1:7771"); err != nil {
		logs.Panic("service run", "error", err)
	}
}
