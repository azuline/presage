package feed

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ParseFeedsList reads the list of feeds from the `feedsList` file and returns
// a parsed list.
func ParseFeedsList(feedsList string) ([]string, error) {
	ccontents, err := os.ReadFile(feedsList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read feeds list")
	}

	log.Println("Parsed feeds list from file")
	return strings.Fields(string(ccontents)), nil
}
