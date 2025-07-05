package services

import (
	"context"
	"errors"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

var (
	ErrCreateTokens = errors.New("error when create tokenst pair in Signup")
)

func (s *services) Signup(ctx context.Context, userId string) (tokens.TokensPair, error) {
	tknPair, err := s.tokenManager.NewTokensPair(userId)
	if err != nil {
		return tknPair, errors.Join(ErrCreateTokens, err)
	}
	return tknPair, nil
}
