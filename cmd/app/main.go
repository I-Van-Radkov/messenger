package main

import (
	"github.com/I-Van-Radkov/messenger/internal/app"
	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.MustLoad()

	application := app.NewApp(cfg)

	application.HTTPServer.MustRun()
}
