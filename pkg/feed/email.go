package feed

import (
	"context"
	"fmt"

	"github.com/azuline/presage/pkg/email"
	"github.com/azuline/presage/pkg/services"
)

func SendEntry(_ context.Context, _ *services.Services, _ email.EmailAddress, _ Entry) error {
	fmt.Printf("ok lol")
	return nil
}
