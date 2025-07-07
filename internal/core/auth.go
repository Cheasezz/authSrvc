package core

import (
	"context"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(ctx context.Context, userId uuid.UUID) (tokens.TokensPair, error)
}

type AuthRepo interface {
	Signup(ctx context.Context, userId uuid.UUID) error
}
