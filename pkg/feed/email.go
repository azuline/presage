package feed

import (
	"context"
	"fmt"

	"github.com/azuline/presage/pkg/services"
)

func SendEntry(_ context.Context, _ *services.Services, _ string, _ Entry) error {
	fmt.Printf("ok lol")
	return nil
}
