package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrSignUp = errors.New("error when exec query in Signup")
)

func (r *Repo) Signup(ctx context.Context, userId uuid.UUID) error {
	query := fmt.Sprintf(`INSERT INTO %s (id) values ($1)`, usersTable)

	_, err := r.PgxPool.Exec(ctx, query, userId)
	if err != nil {
		return errors.Join(ErrSignUp, err)
	}
	return nil
}
