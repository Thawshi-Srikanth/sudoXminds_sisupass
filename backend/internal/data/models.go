package data

import (
	"database/sql"
	"errors"

	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/store/postgres"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users  repository.UserRepository
	Tokens repository.TokenRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  postgres.NewUserRepository(db),
		Tokens: postgres.NewTokenRepository(db),
	}
}
