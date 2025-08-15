package app

import (
	"sync"

	"sisupass.com/email-client/cmd/api/config"
	"sisupass.com/email-client/internal/jsonlog"
	"sisupass.com/email-client/internal/services"
)

const Version = "1.0.0"

type Application struct {
	Config   *config.Config
	Logger   *jsonlog.Logger
	Services *Services
	WG       sync.WaitGroup
}

type Services struct {
	Mail *services.MailerService
}
