package feed

import (
	"context"
	"fmt"

	"github.com/azuline/presage/pkg/services"
)

func SendEntry(_ context.Context, _ *services.Services, to string, entry Entry) error {
	fmt.Printf("to: %s title: %s\n", to, entry.Title)
	return nil
}
