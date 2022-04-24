package feed

import "context"

type Feed struct {
	URL string
}

// ParseFeedsList reads the list of feeds from the `feedsList` file and returns
// a parsed list.
func ParseFeedsList(ctx context.Context, feedsList string) ([]Feed, error) {
	return []Feed{}, nil
}
