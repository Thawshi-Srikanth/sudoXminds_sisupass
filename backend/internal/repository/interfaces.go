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
	FindOrCreateFromGoogle(googleUser *types.GoogleUser) (*types.User, error)
}

type TokenRepository interface {
	Insert(token *tokens.Token) error
	DeleteAllForUser(scope string, userID uuid.UUID) error
	New(userID uuid.UUID, ttl time.Duration, scope string) (*tokens.Token, error)
}

type WalletRepository interface {
	Create(wallet *types.Wallet) error
	GetByID(walletID string) (*types.Wallet, error)
	GetByUserID(userID string) ([]*types.Wallet, error)
	Update(wallet *types.Wallet) error
	Delete(walletID string) error
	GetBalance(walletID string) (float64, error)
	UpdateBalance(walletID string, amount float64) error
}

type SlotRepository interface {
	Create(slot *types.Slot) error
	GetByID(slotID string) (*types.Slot, error)
	GetByUserID(userID string) ([]*types.Slot, error)
	GetByType(userID, slotType string) ([]*types.Slot, error)
	Update(slot *types.Slot) error
	Delete(slotID string) error
	Search(userID, query string) ([]*types.Slot, error)
}
