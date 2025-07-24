package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/authSrvc/internal/core"
)

var (
	ErrSignUp            = errors.New("error when Exec query in Signup")
	ErrGetSessionById    = errors.New("error when Get in GetSessionById")
	ErrDeleteSessionById = errors.New("error when Exec in DeleteSessionById")
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
		(id,user_id,refresh_token,user_agent,ip,expires_at)
		values ($1,$2,$3,$4,$5)`, sessionsTable)
	_, err = tx.Exec(
		ctx,
		query,
		session.Id,
		session.UserId,
		session.RefreshTokenHash,
		session.UserAgent,
		session.Ip,
		session.ExpriresAt,
	)
	if err != nil {
		return errors.Join(ErrSignUp, err)
	}

	return tx.Commit(ctx)
}

func (r *Repo) GetSessionById(ctx context.Context, sessionId string) (*core.Session, error) {
	var session core.Session

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1 `, sessionsTable)

	err := r.DB.Scanny.Get(ctx, r.DB.Pool, &session, query, sessionId)
	if err != nil {
		return nil, ErrGetSessionById
	}

	return &session, nil
}

func (r *Repo) DeleteSessionById(ctx context.Context, sessionId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", sessionsTable)

	_, err := r.DB.Pool.Exec(ctx, query, sessionId)
	if err != nil {
		return ErrDeleteSessionById
	}

	return nil
}
