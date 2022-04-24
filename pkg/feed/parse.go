package feed

import (
	"context"
	"os"
	"strings"
)

type Feed struct {
	URL string
}

// ParseFeedsList reads the list of feeds from the `feedsList` file and returns
// a parsed list.
func ParseFeedsList(ctx context.Context, feedsList string) ([]Feed, error) {
	bcontents, err := os.ReadFile(feedsList)
	if err != nil {
		return nil, err
	}

	scontents := string(bcontents)
	feeds := convertContentsToFeeds(scontents)
	return feeds, nil
}

func convertContentsToFeeds(contents string) []Feed {
	lines := strings.Fields(contents)

	feeds := make([]Feed, len(lines))
	for i, line := range lines {
		feeds[i] = Feed{URL: line}
	}

	return feeds
}
