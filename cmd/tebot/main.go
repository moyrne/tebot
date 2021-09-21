package main

import (
	"context"
	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/moyrne/tebot/api"
	"github.com/moyrne/tebot/configs"
	"github.com/moyrne/tebot/internal/biz"
	"github.com/moyrne/tebot/internal/data"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/pkg/keepalive"
	"github.com/moyrne/tebot/internal/pkg/logs"
	"github.com/moyrne/tebot/internal/pkg/ratelimit"
	"github.com/moyrne/tebot/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// Json 格式化
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 加载配置文件
	if err := configs.LoadConfig(); err != nil {
		logrus.Panicf("read config error %v\n", err)
	}

	// 连接数据库
	if err := database.ConnectMySQL(); err != nil {
		logrus.Panicf("db connect error %v\n", err)
	}

	// 连接Redis
	if err := database.ConnectRedis(); err != nil {
		logrus.Panicf("redis connect error %v\n", err)
	}

	// 初始化日志
	hook, cl, err := logs.NewFileHook(viper.GetString("LogValue.Filename"))
	if err != nil {
		logrus.Panicf("init logrus error %v\n", err)
	}
	defer cl()
	logrus.AddHook(hook)
	logrus.AddHook(biz.NewDBHook(database.DB, data.NewLogRepo()))

	// 心跳检测
	go keepalive.StartCQHTTP()

	// 初始化限流
	ratelimit.InitRate(ratelimit.NewRedisLimit(database.Redis))

	// 启动自动回复 同步
	biz.SyncReply(context.Background(), data.NewReplyRepo())

	// 启动WEB服务
	gin.SetMode(gin.DebugMode)
	e := gin.New()
	e.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/event"}}), gin.Recovery())
	api.RegisterCQServer(e, service.NewEventServer(data.NewEventRepo()))
	if err := e.Run(viper.GetString("Server.Addr")); err != nil {
		logrus.Panicf("service run error %v\n", err)
	}
}
