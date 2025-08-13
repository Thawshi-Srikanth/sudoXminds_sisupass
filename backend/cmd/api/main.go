package main

import (
	"flag"
	"os"

	"sisupass.com/sisupass/cmd/api/app"
	appconfig "sisupass.com/sisupass/cmd/api/config"
	"sisupass.com/sisupass/cmd/api/server"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/jsonlog"
	"sisupass.com/sisupass/internal/migrate"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/migrations"
)

func main() {
	var runMigrations = flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	if logger == nil {
		panic("Logger is not initialized")
	}

	db, err := cfg.OpenDB()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	if *runMigrations {
		logger.PrintInfo("running database migrations", nil)
		err = migrate.MigrateFS(db, migrations.FS, ".")
		if err != nil {
			logger.PrintFatal(err, nil)
		}
		logger.PrintInfo("database migrations completed", nil)
		return
	}

	models := data.NewModels(db)

	appServies := &app.Services{
		Users:  services.NewUserService(&models),
		Tokens: services.NewTokenService(&models),
	}

	app := &app.Application{
		Config:            cfg,
		Logger:            logger,
		Models:            models,
		Services:          appServies,
		// GoogleOAuthConfig: cfg.InitGoogleOAuth(),
	}

	err = server.Serve(app)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
