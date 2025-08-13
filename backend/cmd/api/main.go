// SisuPass API
//
// This is the SisuPass authentication and user management API.
// It provides endpoints for user registration, authentication, password management, and user activation.
//
//	Title:			SisuPass API
//	Description:	SisuPass authentication and user management API
//	Version:		1.0.0
//	Host:			localhost:4000
//	BasePath:		/api/v1
//	Contact:		SisuPass Team <support@sisupass.com>
//
//	Schemes: http, https
//
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	    description: Type "Bearer" followed by a space and JWT token.
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"flag"
	"os"

	_ "github.com/lib/pq"
	"sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/cmd/api/config"
	"sisupass.com/sisupass/cmd/api/server"
	_ "sisupass.com/sisupass/docs"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/jsonlog"
	"sisupass.com/sisupass/internal/mailer"
	"sisupass.com/sisupass/internal/migrate"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/migrations"
)

func main() {
	var runMigrations = flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	cfg, err := config.LoadConfig()
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

	mailSender := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)

	appServies := &app.Services{
		Users:  services.NewUserService(&models),
		Tokens: services.NewTokenService(&models),
		Mail:   services.NewMailService(mailSender, cfg.FrontendURL),
	}

	app := &app.Application{
		Config:   cfg,
		Logger:   logger,
		Models:   models,
		Services: appServies,
		Mailer:   mailSender,
		// GoogleOAuthConfig: cfg.InitGoogleOAuth(),
	}

	err = server.Serve(app)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
