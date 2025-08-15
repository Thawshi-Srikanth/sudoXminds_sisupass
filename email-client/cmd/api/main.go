package main

import (
	"flag"
	"log"
	"os"

	_ "sisupass.com/email-client/docs" 
	"sisupass.com/email-client/cmd/api/app"
	"sisupass.com/email-client/cmd/api/config"
	"sisupass.com/email-client/cmd/api/server"
	"sisupass.com/email-client/internal/jsonlog"
	"sisupass.com/email-client/internal/services"
)

// @title			SisuPass Email Service API
// @version			1.0
// @description		A standalone email service for SisuPass application providing template-based email sending, bulk operations, and delivery tracking.
// @termsOfService	http://swagger.io/terms/

// @contact.name	SisuPass API Support
// @contact.url		http://www.sisupass.com/support
// @contact.email	support@sisupass.com

// @license.name	MIT
// @license.url		http://www.apache.org/licenses/LICENSE-2.0.html

// @host			localhost:4000
// @BasePath		/api/v1

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		log.Printf("Version:\t%s", cfg.Version)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	mailerService := services.NewMailerService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.Sender,
		5, // 5 workers
	)

	application := &app.Application{
		Config: cfg,
		Logger: logger,
		Services: &app.Services{
			Mail: mailerService,
		},
	}

	err = server.Serve(application)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
