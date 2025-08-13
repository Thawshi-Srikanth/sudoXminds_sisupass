package app

import (
	"sync"

	"sisupass.com/sisupass/cmd/api/config"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/jsonlog"
	"sisupass.com/sisupass/internal/mailer"
	"sisupass.com/sisupass/internal/services"
)

const Version = "1.0.0"

type Application struct {
	Config   *config.Config
	Logger   *jsonlog.Logger
	WG       sync.WaitGroup
	Models   data.Models
	Services *Services
	Mailer   mailer.Mailer
}

type Services struct {
	Users  *services.UserService
	Tokens *services.TokenService
	Mail   *services.MailService
}
