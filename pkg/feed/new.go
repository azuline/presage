package feed

import (
	"context"

	"github.com/azuline/presage/pkg/services"
)

func ReadNewEntries(_ context.Context, _ *services.Services) ([]Entry, error) {
	return []Entry{}, nil
}
