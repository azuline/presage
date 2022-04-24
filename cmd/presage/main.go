package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/azuline/presage/pkg/email"
	"github.com/azuline/presage/pkg/feed"
	"github.com/azuline/presage/pkg/migrate"
	"github.com/azuline/presage/pkg/services"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	// Read CLI flags.
	envFile := *flag.String("env-file", "", "path to env file (defaults to .env)")
	feedsList := *flag.String("feeds-list", "", "path to the RSS feeds")
	sendTo := email.EmailAddress(*flag.String("send-to", "", "email to send to"))
	dryRun := *flag.Bool("dry-run", false, "don't send any emails")
	flag.Parse()

	// Load environment variables.
	if err := loadEnvvars(envFile); err != nil {
		log.Fatalln(err)
	}

	// Read environment variables.
	databasePath := os.Getenv("DATABASE_PATH")
	smtpCreds := email.SMTPCreds{
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
	}

	// Initialize DB & Email services.
	srv, err := services.Initialize(databasePath, smtpCreds)
	if err != nil {
		log.Fatalln(err)
	}

	// Migrate database.
	if err := migrate.Migrate(srv); err != nil {
		log.Fatalln(err)
	}

	// Start tool workflow logic composition.
	feeds, err := feed.ParseFeedsList(ctx, feedsList)
	if err != nil {
		log.Fatalln(err)
	}

	if err := feed.DownloadNewFeedEntries(ctx, srv, feeds); err != nil {
		log.Fatalln(err)
	}

	newEntries, err := feed.ReadNewEntries(ctx, srv)
	if err != nil {
		log.Fatalln(err)
	}

	for _, entry := range newEntries {
		log.Printf("Sending new entry %s\n", entry.Link)
		if dryRun {
			continue
		}

		if err := feed.SendEntry(ctx, srv, sendTo, entry); err != nil {
			log.Fatalln(err)
		}
	}
}

func loadEnvvars(envFile string) error {
	// If envFile is passed in, read from that.
	if envFile != "" {
		log.Printf("Loaded environment from %s\n", envFile)
		if err := godotenv.Load(envFile); err != nil {
			errors.Wrap(err, "failed to load envFile")
		}
		return nil
	}

	// Otherwise, if .env exists, read from that.
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		log.Println("Loaded environment from .env")
		if err := godotenv.Load(); err != nil {
			errors.Wrap(err, "failed to load .env")
		}
		return nil
	}

	// Read nothing, default to the normal envvars.
	log.Println("Did not load environment from file")
	return nil
}
