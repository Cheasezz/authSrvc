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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env variables: %s", err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err.Error())
	}

	env, err := app.NewEnv(cfg)
	if err != nil {
		log.Fatalf("Error create environment: %s", err.Error())
	}
	defer env.Close()

	handlers := handlers.New(env)
	srv := httpsrvr.New(handlers.Init(), cfg.HTTP.Host, cfg.HTTP.Port)
	env.Logger.Info("Auth service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
		env.Logger.Info(s.String(), "signal quit")
	case err = <-srv.Notify():
		env.Logger.Error(err, "auth server error on srv.Notify")
	}

	if err := srv.Shutdown(); err != nil {
		env.Logger.Error(err, "auth server error shutting down")
	}

	env.Logger.Info("Auth service server shutting down")
}
