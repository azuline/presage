package feed

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/azuline/presage/pkg/fixtures"
)

func TestReadNewEntries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	srv := fixtures.Services(t)

	feed := TestFeedSource(t, srv)
	// One extra to test against a missing join condition.
	TestFeedSource(t, srv)

	entry1 := TestFeedEntry(t, srv, feed)
	sent1 := TestSentEmail(t, srv, entry1, "")

	TestFeedEntry(t, srv, feed)
	TestFeedEntry(t, srv, feed)
	entry3 := TestFeedEntry(t, srv, feed)
	TestSentEmail(t, srv, entry3, sent1.ToEmail)

	unsent, err := ReadNewEntries(ctx, srv, sent1.ToEmail)
	require.NoError(t, err)

	require.Len(t, unsent, 2)
}
