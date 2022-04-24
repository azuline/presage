package feed

import (
	"context"
	"fmt"
	"log"

	"github.com/azuline/presage/pkg/services"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

var ErrEntryHasNoLink = errors.New("entry has no link")

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
			err := storeEntry(ctx, srv, feed, parsedItem)
			if errors.Is(err, ErrEntryHasNoLink) {
				log.Printf("Failed to store entry b/c no link: %s from %s", parsedItem.Title, url)
				continue
			}
			if err != nil {
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

	res, err := srv.DB.NamedExecContext(ctx, `
		INSERT INTO feed_sources (link, title, authors)
		VALUES (:link, :title, :authors)
		ON CONFLICT (link) DO UPDATE SET 
			title = title,
			authors = authors
	`, feed)
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
	link := parsedItem.Link
	if link == "" && len(parsedItem.Links) > 0 {
		link = parsedItem.Links[0]
	} else {
		return ErrEntryHasNoLink
	}

	entry := Entry{
		FeedID:      feed.ID,
		Link:        link,
		PublishedOn: *parsedItem.PublishedParsed,
		Title:       parsedItem.Title,
		Description: parsedItem.Description,
		Content:     parsedItem.Content,
	}

	res, err := srv.DB.NamedExecContext(ctx, `
		INSERT INTO feed_entries
			(id, source_id, link, published_on, title, description, content)
		VALUES 
			(:id, :source_id, :link, :published_on, :title, :description, :content)
		ON CONFLICT IGNORE
	`, feed)
	if err != nil {
		return errors.Wrap(err, "failed upserting feed entry")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to fetch rows affected")
	}
	if rowsAffected > 0 {
		log.Printf("Stored blog post %s - %s\n", feed.Link, entry.Title)
	}

	return nil
}
