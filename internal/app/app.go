package app

import (
	apphttp "github.com/I-Van-Radkov/messenger/internal/app/http"
	"github.com/I-Van-Radkov/messenger/internal/app/http/handlers"
	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/I-Van-Radkov/messenger/internal/repository/mariadb"
	"github.com/I-Van-Radkov/messenger/internal/services/auth"
)

type App struct {
	HTTPServer *apphttp.Server
}

func NewApp(cfg *config.Config) *App {
	db := mariadb.MustConnect(cfg.DB)

	authService := auth.NewService(
		mariadb.NewUserRepository(db),
		cfg.Auth,
	)

	authHttpHandlers := handlers.NewAuthHandlers(authService)

	httpRouter := apphttp.NewRouterGin(authHttpHandlers)

	httpServer, err := apphttp.NewServer(cfg.HTTP.Port, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout, httpRouter)
	if err != nil {
		panic(err)
	}

	return &App{
		HTTPServer: httpServer,
	}
}
