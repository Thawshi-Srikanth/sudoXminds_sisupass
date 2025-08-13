package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"sisupass.com/sisupass/cmd/api/config"
	"sisupass.com/sisupass/internal/migrate"
	"sisupass.com/sisupass/migrations"
)

func main() {
	var command = flag.String("cmd", "up", "Migration command: up, down, status")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := cfg.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch *command {
	case "up":
		fmt.Println("Running migrations...")
		err = migrate.MigrateFS(db, migrations.FS, ".")
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "down":
		fmt.Println("Rolling back migration...")
		err = migrate.Down(db, ".")
		if err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Rollback completed successfully!")

	case "status":
		fmt.Println("Migration status:")
		err = migrate.Status(db, ".")
		if err != nil {
			log.Fatal("Status check failed:", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, status")
		os.Exit(1)
	}
}
