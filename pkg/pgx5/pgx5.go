package pgx5

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	connAttempts = 10
)

var (
	ErrParseCfg         = errors.New("errpo when parse config in pgx5.New")
	ErrConnAttemptsOver = errors.New("error when connect to db in pgx5.New. connAttempts == 0")
	ErrPing             = errors.New("error when ping db in pgx5.New")
)

type Pgx5 struct {
	Pool *pgxpool.Pool
}

func New(dbUrl string) (*Pgx5, error) {
	var pool *pgxpool.Pool

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, errors.Join(ErrParseCfg, err)
	}

	for connAttempts > 0 {
		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Pgx is trying to connect, attempts left: %d", connAttempts)

		time.Sleep(time.Second)

		connAttempts--
	}

	if err != nil {
		return nil, errors.Join(ErrConnAttemptsOver, err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, errors.Join(ErrPing, err)
	}

	log.Printf("Postgres connected, connAttempts: %d", connAttempts)

	return &Pgx5{Pool: pool}, nil
}

func (p *Pgx5) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
