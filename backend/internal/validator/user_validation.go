package validator

import "sisupass.com/sisupass/internal/types"


func ValidateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(Matches(email, EmailRx), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateCreateUserRequest(v *Validator, username, email, password string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(len(username) >= 3, "username", "must be at least 3 characters long")
	v.Check(len(username) <= 50, "username", "must not be more than 50 characters long")

	ValidateEmail(v, email)
	ValidatePasswordPlaintext(v, password)
}

func ValidateUser(v *Validator, user *types.User) {
	v.Check(user.UserName != "", "username", "must be provided")
	v.Check(len(user.UserName) <= 500, "username", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.HasPlaintext() {
		ValidatePasswordPlaintext(v, user.Password.Plaintext())
	}

	v.Check(user.Password.HasHash(), "password", "password hash is required")

}