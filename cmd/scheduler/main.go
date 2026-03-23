package main

import (
	"log"
	"net/http"

	"github.com/Prizze/TaskScheduler/internal/app"
	"github.com/Prizze/TaskScheduler/internal/config"
)

func main() {
	// Грузим конфиг
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	// Запускаем сервер
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
	defer app.Close()

	if err := app.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run app: %v", err)
	}
}
