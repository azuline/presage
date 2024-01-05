package feed

import (
	"context"
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/azuline/presage/pkg/services"
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
			log.Printf("Failed to parse feed %s: %s", url, err)
			continue
		}

		feed, err := upsertFeed(ctx, srv, url, parsedFeed)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to upsert feed %s", url))
		}

		for _, parsedItem := range parsedFeed.Items {
			_, err := storeEntry(ctx, srv, feed, parsedItem)
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

	_, err := srv.DB.NamedExecContext(ctx, `
		INSERT INTO feed_sources (link, title, authors)
		VALUES (:link, :title, :authors)
		ON CONFLICT (link) DO UPDATE SET 
			title = :title,
			authors = :authors
	`, feed)
	if err != nil {
		return feed, err
	}

	err = srv.DB.Get(&feed, `
		SELECT id, link, title, authors
		FROM feed_sources
		WHERE link = ?
	`, url)
	if err != nil {
		return feed, err
	}

	log.Printf("Downloaded and updated feed details: %s [%s]", feed.Title, feed.Link)
	return feed, nil
}

func storeEntry(
	ctx context.Context,
	srv *services.Services,
	feed Feed,
	parsedItem *gofeed.Item,
) (int, error) {
	link := parsedItem.Link
	if link == "" {
		if len(parsedItem.Links) > 0 {
			link = parsedItem.Links[0]
		} else {
			return 0, ErrEntryHasNoLink
		}
	}

	publishedOn := "Unknown"
	if parsedItem.PublishedParsed != nil {
		publishedOn = parsedItem.PublishedParsed.Format("02 Jan 06")
	}

	entry := Entry{
		FeedID:      feed.ID,
		Link:        link,
		PublishedOn: publishedOn,
		Title:       parsedItem.Title,
		Description: parsedItem.Description,
		Content:     parsedItem.Content,
	}

	res, err := srv.DB.NamedExecContext(ctx, `
		INSERT INTO feed_entries
			(source_id, link, published_on, title, description, content)
		VALUES 
			(:source_id, :link, :published_on, :title, :description, :content)
		ON CONFLICT (link) DO NOTHING
	`, entry)
	if err != nil {
		return 0, errors.Wrap(err, "failed upserting feed entry")
	}

	entryID, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch last insert ID")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch rows affected")
	}
	if rowsAffected > 0 {
		log.Printf("Stored feed entry %d %s - %s\n", entryID, feed.Link, entry.Title)
	}

	return int(entryID), nil
}
