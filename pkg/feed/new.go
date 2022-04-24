package feed

import (
	"context"

	"github.com/azuline/presage/pkg/services"
)

// An entry is an entry in a feed. Probably an article.
type Entry struct {
	Source Feed
	Author string
	Link   string
	Body   string
}

func ReadNewEntries(_ context.Context, _ *services.Services) ([]Entry, error) {
	return []Entry{}, nil
}
