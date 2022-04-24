package feed

import (
	"context"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"

	"github.com/azuline/presage/pkg/fixtures"
)

func TestUpsertFeed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srv := fixtures.Services(t)
	url := "http://feed.one/rss.xml"

	var firstFeed Feed
	t.Run("first feed insert", func(t *testing.T) {
		parsedFeed := &gofeed.Feed{
			Title: "Feed Title",
			Authors: []*gofeed.Person{
				{
					Name:  "Person 1",
					Email: "person1@ema.il",
				},
				{
					Name:  "Person 2",
					Email: "person2@ema.il",
				},
			},
		}

		var err error
		firstFeed, err = upsertFeed(ctx, srv, url, parsedFeed)
		require.NoError(t, err)
		require.Equal(t, 1, firstFeed.ID)
		require.Equal(t, url, firstFeed.Link)
		require.Equal(t, parsedFeed.Title, firstFeed.Title)
		require.Equal(t, "Person 1 & Person 2", firstFeed.Authors)
	})

	t.Run("same feed upsert conflict", func(t *testing.T) {
		parsedFeed := &gofeed.Feed{
			Title: "New Title",
			Authors: []*gofeed.Person{
				{
					Name:  "Person 3",
					Email: "person1@ema.il",
				},
			},
		}
		var err error
		firstFeed, err = upsertFeed(ctx, srv, url, parsedFeed)
		require.NoError(t, err)
		require.Equal(t, 1, firstFeed.ID)
		require.Equal(t, url, firstFeed.Link)
		require.Equal(t, parsedFeed.Title, firstFeed.Title)
		require.Equal(t, "Person 3", firstFeed.Authors)
	})
}

func TestStoreEntry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srv := fixtures.Services(t)
	feed := TestFeedSource(t, srv)

	now := time.Now()
	firstParsedItem := &gofeed.Item{
		Link:            "https://feed.one/article-one",
		PublishedParsed: &now,
		Title:           "Title",
		Description:     "Interesting description",
		Content:         "3000 words go here",
	}

	firstEntry := Entry{}
	t.Run("insert first entry", func(t *testing.T) {
		entryID, err := storeEntry(ctx, srv, feed, firstParsedItem)
		require.NoError(t, err)

		err = srv.DB.Get(&firstEntry, `
			SELECT id, source_id, link, published_on, title, description, content
			FROM feed_entries
			WHERE id = ?
		`, entryID)
		require.NoError(t, err)

		require.Equal(t, feed.ID, firstEntry.FeedID)
		require.Equal(t, firstParsedItem.Link, firstEntry.Link)
		require.True(t, firstEntry.PublishedOn.Equal(*firstParsedItem.PublishedParsed))
		require.Equal(t, firstParsedItem.Title, firstEntry.Title)
		require.Equal(t, firstParsedItem.Description, firstEntry.Description)
		require.Equal(t, firstParsedItem.Content, firstEntry.Content)
	})

	t.Run("insert duplicate entry", func(t *testing.T) {
		entryID, err := storeEntry(ctx, srv, feed, firstParsedItem)
		require.NoError(t, err)
		require.Equal(t, firstEntry.ID, entryID)
	})
}
