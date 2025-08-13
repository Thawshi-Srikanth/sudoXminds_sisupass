package repository

import (
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/tokens"
	"sisupass.com/sisupass/internal/types"
)

type UserRepository interface {
	GetByEmail(email string) (*types.User, error)
	GetForToken(tokenScope, tokenPlainText string) (*types.User, error)
	GetByID(id string) (*types.User, error)
	Create(user *types.User) error
	Update(user *types.User) error
	Delete(id string) error
}

type TokenRepository interface {
	Insert(token *tokens.Token) error
	DeleteAllForUser(scope string, userID uuid.UUID) error
	New(userID uuid.UUID, ttl time.Duration, scope string) (*tokens.Token, error)
}
