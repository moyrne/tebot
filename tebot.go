package main

import (
	"github.com/moyrne/tebot/internal/service/api"
	"os"
)

func main() {
	r := api.NewRouter()
	if err := r.Run("127.0.0.1:7771"); err != nil {
		os.Exit(1)
		return
	}
}
