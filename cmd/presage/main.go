package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/azuline/presage/pkg/email"
	"github.com/azuline/presage/pkg/feed"
	"github.com/azuline/presage/pkg/services"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	// Read CLI flags.
	feedsList := *flag.String("feeds-list", "", "path to the RSS feeds")
	sendTo := email.EmailAddress(*flag.String("send-to", "", "email to send to"))
	dryRun := *flag.Bool("dry-run", false, "don't send any emails")
	flag.Parse()

	// Read environment variables.
	sendgridKey := os.Getenv("SENDGRID_KEY")
	databaseURI := os.Getenv("DATABASE_URI")

	srv, err := services.Initialize(sendgridKey, databaseURI)
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
