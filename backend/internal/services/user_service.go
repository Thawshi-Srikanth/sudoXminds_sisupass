package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/tokens"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

type UserService struct {
	models *data.Models
}

func NewUserService(models *data.Models) *UserService {
	return &UserService{models: models}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	UserType  string    `json:"user_type"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {

	v := validator.New()
	validator.ValidateCreateUserRequest(v, req.Username, req.Email, req.Password)

	if !v.Valid() {
		return nil, validator.ValidationError{Errors: v.Errors}
	}

	_, err := s.models.Users.GetByEmail(req.Email)
	if err == nil {
		return nil, types.ErrDuplicateEmail
	} else if err != types.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	user := &types.User{
		ID:                    uuid.New(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		UserName:              req.Username,
		Email:                 req.Email,
		UserType:              "student",
		Activated:             false,
		AuthenticationMethods: "password",
		Version:               1,
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userValidator := validator.New()
	validator.ValidateUser(userValidator, user)

	if !userValidator.Valid() {
		return nil, validator.ValidationError{Errors: userValidator.Errors}
	}

	err = s.models.Users.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.UserName,
		Email:     user.Email,
		UserType:  user.UserType,
		Activated: user.Activated,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {

	v := validator.New()
	validator.ValidateEmail(v, email)
	if !v.Valid() {
		return nil, validator.ValidationError{Errors: v.Errors}
	}

	user, err := s.models.Users.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.UserName,
		Email:     user.Email,
		UserType:  user.UserType,
		Activated: user.Activated,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return validator.ValidationError{
			Errors: map[string]string{
				"user_id": "must be a valid UUID",
			},
		}
	}

	_, err := s.models.Users.GetByID(userID)
	if err != nil {
		return err
	}

	err = s.models.Users.Delete(userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

type ActivateUserRequest struct {
	TokenPlaintext string `json:"token"`
}

func (s *UserService) ActivateUser(ctx context.Context, req ActivateUserRequest) (*UserResponse, error) {
	v := validator.New()
	validator.ValidateTokenPlaintext(v, req.TokenPlaintext)

	if !v.Valid() {
		return nil, validator.ValidationError{Errors: v.Errors}
	}

	user, err := s.models.Users.GetForToken(tokens.ScopeActivation, req.TokenPlaintext)
	if err != nil {
		if err == types.ErrRecordNotFound {
			return nil, validator.ValidationError{
				Errors: map[string]string{
					"token": "invalid or expired activation token",
				},
			}
		}
		return nil, fmt.Errorf("failed to get user for token: %w", err)
	}

	user.Activated = true

	err = s.models.Users.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	err = s.models.Tokens.DeleteAllForUser(tokens.ScopeActivation, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete activation tokens: %w", err)
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.UserName,
		Email:     user.Email,
		UserType:  user.UserType,
		Activated: user.Activated,
		CreatedAt: user.CreatedAt,
	}, nil
}

type UpdatePasswordRequest struct {
	Password       string `json:"password"`
	TokenPlaintext string `json:"token"`
}

func (s *UserService) UpdateUserPassword(ctx context.Context, req UpdatePasswordRequest) error {
	v := validator.New()
	validator.ValidatePasswordPlaintext(v, req.Password)
	validator.ValidateTokenPlaintext(v, req.TokenPlaintext)

	if !v.Valid() {
		return validator.ValidationError{Errors: v.Errors}
	}

	user, err := s.models.Users.GetForToken(tokens.ScopePasswordReset, req.TokenPlaintext)
	if err != nil {
		if err == types.ErrRecordNotFound {
			return validator.ValidationError{
				Errors: map[string]string{
					"token": "invalid or expired password reset token",
				},
			}
		}
		return fmt.Errorf("failed to get user for token: %w", err)
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	err = s.models.Users.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	err = s.models.Tokens.DeleteAllForUser(tokens.ScopePasswordReset, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete password reset tokens: %w", err)
	}

	return nil
}
