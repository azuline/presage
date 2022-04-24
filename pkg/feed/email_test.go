package feed

import (
	"context"
	"testing"

	"github.com/azuline/presage/pkg/fixtures"
	"github.com/stretchr/testify/require"
)

func TestRecordSentEmail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srv := fixtures.Services(t)

	feed := TestFeedSource(t, srv)
	rawEntry := TestFeedEntry(t, srv, feed)

	entry := EntryWithSourceTitle{
		Entry:       rawEntry,
		SourceTitle: feed.Title,
	}

	err := recordSentEntry(ctx, srv, "you@ema.il", entry)
	require.NoError(t, err)

	rows := srv.DB.QueryRowx(`
		SELECT id
		FROM feed_sent_emails
		WHERE entry_id = ?
	`, entry.ID)
	require.NoError(t, rows.Err())

	var id int
	err = rows.Scan(&id)
	require.NoError(t, err)
	require.Equal(t, entry.ID, id)
}
