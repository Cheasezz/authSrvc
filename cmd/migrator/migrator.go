package main

import (
	"errors"
	"log"
	"time"

	"github.com/Cheasezz/authSrvc/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var (
	attempts = 10
	m        *migrate.Migrate
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env variables: %s", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	for attempts > 0 {
		m, err = migrate.New("file://"+cfg.PG.Schema_Url, "pgx5://"+cfg.PG.URL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d \n", attempts)
		log.Println("err from migrate new error : ", err)
		time.Sleep(time.Second)
		attempts--
	}

	if err != nil {
		log.Fatal("Error when migrate ", err)
	}
	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate: no change")
		return
	}

	log.Println("Migrate: up success")
}
