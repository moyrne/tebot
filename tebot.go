package main

import (
	"github.com/moyrne/tebot/configs"
	"github.com/moyrne/tebot/internal/database"
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tebot/internal/service/api"
	"log"
	"os"
)

func main() {
	if err := configs.LoadConfig(); err != nil {
		log.Fatalln("read config error", err)
	}
	writer, err := logs.FileWriter()
	if err != nil {
		log.Fatalln("new file writer error", err)
	}
	logs.Init(writer)
	if err := database.ConnectPG(); err != nil {
		log.Fatalln("db connect error", err)
	}
	r := api.NewRouter()
	if err := r.Run("127.0.0.1:7771"); err != nil {
		logs.Error("service run", "error", err)
		os.Exit(1)
	}
}
