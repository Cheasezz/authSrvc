package repo

import "github.com/jackc/pgx/v5/pgxpool"

const (
	usersTable = "users"
)

type Repo struct {
	PgxPool *pgxpool.Pool
}

func New(pg *pgxpool.Pool) *Repo {
	return &Repo{
		PgxPool: pg,
	}
}
