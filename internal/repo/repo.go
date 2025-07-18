package repo

import "github.com/Cheasezz/authSrvc/pkg/pgx5"

const (
	usersTable    = "users"
	sessionsTable = "users_sessions"
)

type Repo struct {
	DB *pgx5.Pgx5
}

func New(pgx *pgx5.Pgx5) *Repo {
	return &Repo{
		DB: pgx,
	}
}
