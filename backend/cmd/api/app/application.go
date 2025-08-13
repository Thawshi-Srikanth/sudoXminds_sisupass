package app

import (
	"sync"

	"sisupass.com/sisupass/cmd/api/config"
	"sisupass.com/sisupass/internal/jsonlog"
)

const Version = "1.0.0"

type Application struct {
	Config *config.Config
	Logger *jsonlog.Logger
	WG     sync.WaitGroup
}
