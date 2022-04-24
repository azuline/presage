package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/azuline/presage/pkg/email"
	"github.com/azuline/presage/pkg/feed"
	"github.com/azuline/presage/pkg/services"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
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
	err := godotenv.Load(envFile, ".env")
	if err != nil {
		log.Fatalln(err)
	}

	// Read environment variables.
	databaseURI := os.Getenv("DATABASE_URI")
	smtpCreds := email.SMTPCreds{
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}

	srv, err := services.Initialize(databaseURI, smtpCreds)
	if err != nil {
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
