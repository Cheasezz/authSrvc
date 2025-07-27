package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/authSrvc/internal/core"
)

var (
	ErrSignUp            = errors.New("error when Exec query in CreateSession")
	ErrGetSessionById    = errors.New("error when Get in GetSessionById")
	ErrDeleteSessionById = errors.New("error when Exec in DeleteSessionById")
)

func (r *Repo) CreateSession(ctx context.Context, session *core.Session) error {

	query := fmt.Sprintf(`INSERT INTO %s
		(id,user_id,refresh_token_hash,user_agent,ip,expires_at)
		values ($1,$2,$3,$4,$5,$6)`, sessionsTable)

	_, err := r.DB.Pool.Exec(
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

	return nil
}

func (r *Repo) GetSessionById(ctx context.Context, sessionId string) (*core.Session, error) {
	var session core.Session

	query := fmt.Sprintf(`SELECT 
	id,
	user_id,
	refresh_token_hash,
	user_agent,
	ip::text,
	created_at,
	expires_at FROM %s WHERE id = $1 `, sessionsTable)

	err := r.DB.Scany.Get(ctx, r.DB.Pool, &session, query, sessionId)
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
