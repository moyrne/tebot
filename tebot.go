package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/moyrne/tebot/internal/analyze"

	"github.com/moyrne/tebot/configs"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/service/api"
	"log"
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
	if err := database.ConnectPG(); err != nil {
		logs.Panic("db connect", "error", err)
	}
	analyze.SyncReply(context.Background())
	if err := database.ConnectRedis(); err != nil {
		logs.Panic("redis connect", "error", err)
	}
	r := api.NewRouter()
	if err := r.Run("127.0.0.1:7771"); err != nil {
		logs.Panic("service run", "error", err)
	}
}
