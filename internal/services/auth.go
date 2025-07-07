package services

import (
	"context"
	"errors"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
	"github.com/google/uuid"
)

var (
	ErrCreateTokens = errors.New("error when create tokenst pair in Signup")
	ErrRepoSignup   = errors.New("error when repo.Signup in Signup")
)

func (s *services) Signup(ctx context.Context, userId uuid.UUID) (tokens.TokensPair, error) {
	tknPair, err := s.tokenManager.NewTokensPair(userId.String())
	if err != nil {
		return tknPair, errors.Join(ErrCreateTokens, err)
	}
	err = s.repo.Signup(ctx, userId)
	if err != nil {
		return tknPair, errors.Join(ErrRepoSignup, err)
	}
	return tknPair, nil
}
