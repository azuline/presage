package feed

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/azuline/presage/pkg/psrand"
	"github.com/azuline/presage/pkg/services"
)

func TestFeed(t *testing.T, srv *services.Services) Feed {
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
