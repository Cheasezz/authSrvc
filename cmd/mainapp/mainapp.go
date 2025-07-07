package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheasezz/authSrvc/config"
	"github.com/Cheasezz/authSrvc/internal/app"
	"github.com/Cheasezz/authSrvc/internal/handlers"
	"github.com/Cheasezz/authSrvc/pkg/httpsrvr"
	"github.com/joho/godotenv"
)

// @title Auth server API
// @version 1.0
// @description API Server for Auth
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env variables: %s", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	env, err := app.NewEnv(cfg)
	if err != nil {
		log.Fatalf("Error create environment: %s", err)
	}
	defer env.Close()

	handlers := handlers.New(env)
	srv := httpsrvr.New(handlers.Init(cfg.App.DevMod), cfg.HTTP.Host, cfg.HTTP.Port)
	env.Logger.Info("Auth server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
		env.Logger.Info("signal quit: %s", s)
	case err = <-srv.Notify():
		env.Logger.Error("auth server error on srv.Notify: %s", err)
	}

	if err := srv.Shutdown(); err != nil {
		env.Logger.Error("auth server error shutting down: %s", err)
	}

	env.Logger.Info("Auth server shutting down")
}
