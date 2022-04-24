package feed

import (
	"context"

	"github.com/azuline/presage/pkg/services"
)

// DownloadNewFeedEntries downloads RSS feeds from the passed in list and
// stores all new articles into the SQLite database.
func DownloadNewFeedEntries(_ context.Context, _ *services.Services, _ []Feed) error {
	return nil
}
