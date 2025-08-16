package config

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Version string
	Port    int
	Env     string
	DB      struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	OAuth struct {
		GoogleClientID     string
		GoogleClientSecret string
		RedirectURI        string
	}
	FrontendURL string
	CORS        struct {
		TrustedOrigins []string
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.Version = "1.0.0"

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production|testing)")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.SMTP.Host, "smtp-host", os.Getenv("SMTPHOST"), "SMTP host")
	flag.StringVar(&cfg.FrontendURL, "frontend-url", os.Getenv("FRONTEND_URL"), "Frontend URL")

	// Handle SMTP port with environment variable parsing
	envSMTPPort := os.Getenv("SMTPPORT")
	if envSMTPPort == "" {
		envSMTPPort = "587"
		fmt.Println("SMTPPORT is not set. Defaulting to 587")
	}

	intSMTPPort, err := strconv.Atoi(envSMTPPort)
	if err != nil {
		fmt.Println("SMTPPORT is not a number. Defaulting to 587")
		intSMTPPort = 587
	}

	flag.IntVar(&cfg.SMTP.Port, "smtp-port", intSMTPPort, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", os.Getenv("SMTPUSERNAME"), "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", os.Getenv("SMTPPASS"), "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", os.Getenv("SMTPSENDER"), "SMTP sender")

	// OAuth configuration
	flag.StringVar(&cfg.OAuth.GoogleClientID, "oauth-google-client-id", os.Getenv("GOOGLE_CLIENT_ID"), "Google OAuth Client ID")
	flag.StringVar(&cfg.OAuth.GoogleClientSecret, "oauth-google-client-secret", os.Getenv("GOOGLE_CLIENT_SECRET"), "Google OAuth Client Secret")
	flag.StringVar(&cfg.OAuth.RedirectURI, "oauth-redirect-url", os.Getenv("GOOGLE_REDIRECT_URI"), "OAuth Redirect URL")

	// CORS configuration
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.CORS.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	// Set default CORS origins
	cfg.CORS.TrustedOrigins = append(cfg.CORS.TrustedOrigins, "http://localhost:5173", "http://localhost:3000")

	return &cfg, nil
}

func (cfg *Config) OpenDB() (*sql.DB, error) {
	dsn := cfg.DB.DSN
	if !strings.Contains(dsn, "sslmode=") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=disable"
		} else {
			dsn += "?sslmode=disable"
		}
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (cfg *Config) InitGoogleOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.OAuth.GoogleClientID,
		ClientSecret: cfg.OAuth.GoogleClientSecret,
		RedirectURL:  cfg.OAuth.RedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
