package services

import (
	"context"
	"fmt"

	"sisupass.com/sisupass/internal/mailer"
)

type MailService struct {
	mailer      mailer.Mailer
	frontendURL string
}

func NewMailService(mailer mailer.Mailer, frontendURL string) *MailService {
	return &MailService{
		mailer:      mailer,
		frontendURL: frontendURL,
	}
}

type WelcomeEmailData struct {
	Name            string
	ActivationToken string
	ActivationURL   string
}

type PasswordResetEmailData struct {
	Name       string
	ResetToken string
	ResetURL   string
}

func (s *MailService) SendWelcomeEmail(ctx context.Context, recipient, name, activationToken string) error {
	data := WelcomeEmailData{
		Name:            name,
		ActivationToken: activationToken,
		ActivationURL:   fmt.Sprintf("%s/activate?token=%s", s.frontendURL, activationToken),
	}

	return s.mailer.Send(recipient, "user_welcome", data)
}

func (s *MailService) SendPasswordResetEmail(ctx context.Context, recipient, name, resetToken string) error {
	data := PasswordResetEmailData{	

		Name:       name,
		ResetToken: resetToken,
		ResetURL:   fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, resetToken),
	}

	return s.mailer.Send(recipient, "password_reset", data)
}
