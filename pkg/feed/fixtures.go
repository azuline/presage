package feed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/azuline/presage/pkg/psrand"
	"github.com/azuline/presage/pkg/services"
)

func TestFeedSource(t *testing.T, srv *services.Services) Feed {
	feed := Feed{
		Link:    psrand.String(12),
		Title:   psrand.String(12),
		Authors: psrand.String(12),
	}
	res, err := srv.DB.NamedExec(`
		INSERT INTO feed_sources (link, title, authors) 
		VALUES (:link, :title, :authors)
	`, feed)
	require.NoError(t, err)

	lastInsertID, err := res.LastInsertId()
	require.NoError(t, err)

	feed.ID = int(lastInsertID)
	return feed
}

func TestFeedEntry(t *testing.T, srv *services.Services, feed Feed) Entry {
	entry := Entry{
		FeedID:      feed.ID,
		Link:        "https://feed.one/article-one",
		PublishedOn: time.Now(),
		Title:       "Title",
		Description: "Interesting description",
		Content:     "3000 words go here",
	}
	res, err := srv.DB.NamedExec(`
		INSERT INTO feed_entries 
			(source_id, link, published_on, title, description, content)
		VALUES 
			(:source_id, :link, :published_on, :title, :description, :content)
	`, entry)
	require.NoError(t, err)

	lastInsertID, err := res.LastInsertId()
	require.NoError(t, err)

	entry.ID = int(lastInsertID)
	return entry
}
