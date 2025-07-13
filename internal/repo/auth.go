package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/authSrvc/internal/core"
)

var (
	ErrSignUp = errors.New("error when exec query in Signup")
)

func (r *Repo) Signup(ctx context.Context, session core.Session) error {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(`INSERT INTO %s (id) values ($1)`, usersTable)
	_, err = tx.Exec(ctx, query, session.UserId)
	if err != nil {
		return errors.Join(ErrSignUp, err)
	}

	query = fmt.Sprintf(`INSERT INTO %s
		(user_id,refresh_token,user_agent,ip,expires_at)
		values ($1,$2,$3,$4,$5)`, sessionsTable)
	_, err = tx.Exec(
		ctx,
		query,
		session.UserId,
		session.RefreshToken,
		session.UserAgent,
		session.Ip,
		session.ExpriresAt,
	)
	if err != nil {
		return errors.Join(ErrSignUp, err)
	}

	return tx.Commit(ctx)
}
