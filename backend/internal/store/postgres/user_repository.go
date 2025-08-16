package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/types"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*types.User, error) {
	query := `SELECT id, created_at, user_name, email, password_hash, user_type, activated, version FROM users WHERE email = $1`

	var user types.User
	var passwordHash []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UserName,
		&user.Email,
		&passwordHash,
		&user.UserType,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, types.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	user.Password = types.NewPasswordFromHash(passwordHash)

	return &user, nil
}
func (r *UserRepository) GetByID(id string) (*types.User, error) {
	// Implementation
	query := `SELECT id, created_at, updated_at, user_name, email, password_hash, user_type, activated, authentication_methods, version 
              FROM users 
              WHERE id = $1`

	var user types.User
	var passwordHash []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.UserName,
		&user.Email,
		&passwordHash,
		&user.UserType,
		&user.Activated,
		&user.AuthenticationMethods,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, types.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	user.Password = types.NewPasswordFromHash(passwordHash)

	return &user, nil
}

func (r *UserRepository) Create(user *types.User) error {
	// Implementation
	query := `INSERT INTO users (user_name, email, password_hash, user_type, authentication_methods, activated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, version`

	args := []any{user.UserName, user.Email, user.Password.Hash(), user.UserType, user.AuthenticationMethods, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return types.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) Update(user *types.User) error {
	query := `UPDATE users SET user_name = $1, email = $2, password_hash = $3, user_type = $4, activated = $5, version = version + 1 WHERE id = $6 AND version = $7 RETURNING version`

	args := []any{
		user.UserName,
		user.Email,
		user.Password.Hash(),
		user.UserType,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return types.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return types.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (r *UserRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM user_profiles WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return types.ErrRecordNotFound
	}

	return tx.Commit()
}

func (r *UserRepository) GetForToken(tokenScope, tokenPlainText string) (*types.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
        SELECT users.id, users.created_at, users.updated_at, users.user_name, 
               users.email, users.password_hash, users.user_type, users.activated, 
               users.authentication_methods, users.version
        FROM users
        INNER JOIN tokens ON users.id = tokens.user_id
        WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user types.User
	var passwordHash []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.UserName,
		&user.Email,
		&passwordHash,
		&user.UserType,
		&user.Activated,
		&user.AuthenticationMethods,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, types.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	user.Password = types.NewPasswordFromHash(passwordHash)
	return &user, nil
}

func (r *UserRepository) FindOrCreateFromGoogle(googleUser *types.GoogleUser) (*types.User, error) {
	// First, try to find existing user by email
	existingUser, err := r.GetByEmail(googleUser.Email)
	if err == nil {
		// User exists, update their information if needed
		existingUser.Activated = true // Google users are automatically activated
		err = r.Update(existingUser)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	} else if !errors.Is(err, types.ErrRecordNotFound) {
		return nil, err
	}

	// User doesn't exist, create a new one
	user := &types.User{
		UserName:              googleUser.Name,
		Email:                 googleUser.Email,
		UserType:              "user",
		Activated:             true, // Google users are automatically activated
		AuthenticationMethods: "google",
		Version:               1,
	}

	// Set a placeholder password (Google users don't need password login)
	err = user.Password.Set("oauth_google_user_no_password")
	if err != nil {
		return nil, err
	}

	err = r.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
