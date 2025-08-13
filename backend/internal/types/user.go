package types

import (
	"errors"
	"time"

	"github.com/google/uuid" // Use ONE UUID library consistently
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
)

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
	ID                    uuid.UUID `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	UserName              string    `json:"username" db:"user_name"`
	Email                 string    `json:"email" db:"email"`
	Password              Password  `json:"-"`
	UserType              string    `json:"user_type" db:"user_type"`
	Activated             bool      `json:"activated" db:"activated"`
	HasProfileCreated     bool      `json:"has_profile_created" db:"has_profile_created"`
	AuthenticationMethods string    `json:"authentication_methods" db:"authentication_methods"`
	Version               int       `json:"version" db:"version"`
}

type UserProfile struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	MobileNumber *string    `json:"mobile_number,omitempty" db:"mobile_number"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Avatar       *string    `json:"avatar,omitempty" db:"avatar"`
	Bio          *string    `json:"bio,omitempty" db:"bio"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Version      int        `json:"version" db:"version"`
}

type AddressType string

const (
	AddressTypeHome     AddressType = "home"
	AddressTypeWork     AddressType = "work"
	AddressTypeShipping AddressType = "shipping"
	AddressTypeBilling  AddressType = "billing"
)

type Address struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	UserID       uuid.UUID   `json:"user_id" db:"user_id"`
	Type         AddressType `json:"type" db:"type"`
	Label        string      `json:"label" db:"label"`
	FirstName    string      `json:"first_name" db:"first_name"`
	LastName     string      `json:"last_name" db:"last_name"`
	Company      *string     `json:"company,omitempty" db:"company"`
	AddressLine1 string      `json:"address_line_1" db:"address_line_1"`
	AddressLine2 *string     `json:"address_line_2,omitempty" db:"address_line_2"`
	City         string      `json:"city" db:"city"`
	State        string      `json:"state" db:"state"`
	PostalCode   string      `json:"postal_code" db:"postal_code"`
	Country      string      `json:"country" db:"country"`
	PhoneNumber  *string     `json:"phone_number,omitempty" db:"phone_number"`
	IsDefault    bool        `json:"is_default" db:"is_default"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	Version      int         `json:"version" db:"version"`
}
