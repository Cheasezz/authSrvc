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
		log.Printf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	env, err := app.NewEnv(cfg)
	if err != nil {
		log.Fatalf("cannot create environment: %s", err)
	}
	defer env.Close()

	handlers := handlers.New(env)
	srv := httpsrvr.New(handlers.Init(), cfg.HTTP.Host, cfg.HTTP.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
	case err = <-srv.Notify():
	}

	if err := srv.Shutdown(); err != nil {
	}

}
