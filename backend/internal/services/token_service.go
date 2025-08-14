package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/tokens"
)

type TokenService struct {
	models *data.Models
}

func NewTokenService(models *data.Models) *TokenService {
	return &TokenService{models: models}
}

func (s *TokenService) CreateActivationToken(ctx context.Context, userID uuid.UUID) (*tokens.Token, error) {
	err := s.models.Tokens.DeleteAllForUser(tokens.ScopeActivation, userID)
	if err != nil {
		return nil, err
	}

	token, err := s.models.Tokens.New(userID, 3*24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *TokenService) CreateAuthenticationToken(ctx context.Context, userID uuid.UUID) (*tokens.Token, error) {
	token, err := s.models.Tokens.New(userID, 24*time.Hour, tokens.ScopeAuthentication)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *TokenService) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID) (*tokens.Token, error) {
	err := s.models.Tokens.DeleteAllForUser(tokens.ScopePasswordReset, userID)
	if err != nil {
		return nil, err
	}

	token, err := s.models.Tokens.New(userID, 45*time.Minute, tokens.ScopePasswordReset)
	if err != nil {
		return nil, err
	}

	return token, nil
}
