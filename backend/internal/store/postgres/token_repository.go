package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/tokens"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) New(userID uuid.UUID, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *TokenRepository) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TokenRepository) DeleteAllForUser(scope string, userID uuid.UUID) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, scope, userID)
	return err
}
