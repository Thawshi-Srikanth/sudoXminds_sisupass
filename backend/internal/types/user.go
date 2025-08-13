package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
)

// GoogleUser represents user data from Google OAuth
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (p *Password) Hash() []byte {
	return p.hash
}

func (p *Password) SetHash(hash []byte) {
	p.hash = hash
	p.plaintext = nil
}

func (p *Password) Plaintext() string {
	if p.plaintext == nil {
		return ""
	}
	return *p.plaintext
}

func (p *Password) HasPlaintext() bool {
	return p.plaintext != nil
}

func (p *Password) HasHash() bool {
	return len(p.hash) > 0
}

func NewPasswordFromHash(hash []byte) Password {
	return Password{hash: hash}
}

type User struct {
	ID                    uuid.UUID `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	UserName              string    `json:"username"`
	Email                 string    `json:"email"`
	Password              Password  `json:"-"`
	UserType              string    `json:"user_type"`
	Activated             bool      `json:"activated"`
	AuthenticationMethods string    `json:"authentication_methods"`
	Version               int       `json:"version"`
}

type UserProfile struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	MobileNumber *string    `json:"mobile_number,omitempty"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	Bio          *string    `json:"bio,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Version      int        `json:"version"`
}
