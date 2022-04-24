package feed

import (
	"context"
	"fmt"
	"log"

	"github.com/azuline/presage/pkg/services"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// DownloadNewFeedEntries downloads RSS feeds from the passed in list and
// stores all new articles into the SQLite database.
func DownloadNewFeedEntries(
	ctx context.Context,
	srv *services.Services,
	feedURLs []string,
) error {
	parser := gofeed.NewParser()

	for _, url := range feedURLs {
		parsedFeed, err := parser.ParseURLWithContext(url, ctx)
		if err != nil {
			log.Printf("Failed to parse feed %s\n", url)
			continue
		}

		feed, err := upsertFeed(ctx, srv, url, parsedFeed)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to upsert feed %s", url))
		}

		for _, parsedItem := range parsedFeed.Items {
			if err := storeEntry(ctx, srv, feed, parsedItem); err != nil {
				return errors.Wrap(err, fmt.Sprintf(
					"failed to store feed entry %s from %s", parsedItem.Title, url,
				))
			}
		}
	}
	return nil
}

func upsertFeed(
	ctx context.Context,
	srv *services.Services,
	url string,
	parsedFeed *gofeed.Feed,
) (Feed, error) {
	authors := ""
	for _, a := range parsedFeed.Authors {
		if authors == "" {
			authors = a.Name
		} else {
			authors = authors + " & " + a.Name
		}
	}

	feed := Feed{
		ID:      0,
		Link:    url,
		Title:   parsedFeed.Title,
		Authors: authors,
	}

	res, err := srv.DB.ExecContext(ctx, `
		INSERT INTO feed_sources (link, title, authors)
		VALUES (?, ?, ?)
		ON CONFLICT (link) DO UPDATE SET 
			title = title,
			authors = authors
	`, feed.Link, feed.Title, feed.Authors)
	if err != nil {
		return feed, err
	}

	feedID, err := res.LastInsertId()
	if err != nil {
		return feed, err
	}
	feed.ID = int(feedID)

	return feed, nil
}

func storeEntry(
	ctx context.Context,
	srv *services.Services,
	feed Feed,
	parsedItem *gofeed.Item,
) error {
	return nil
}
