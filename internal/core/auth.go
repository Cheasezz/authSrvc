package core

import (
	"context"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type AuthService interface {
	Signup(ctx context.Context, userId string) (tokens.TokensPair, error)
}
