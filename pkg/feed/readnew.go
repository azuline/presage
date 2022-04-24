package feed

import (
	"context"

	"github.com/azuline/presage/pkg/services"
)

// ReadNewEntries returns all feed entries that have not yet been sent out over
// email.
func ReadNewEntries(ctx context.Context, srv *services.Services, notSentTo string) ([]Entry, error) {
	var newEntries []Entry

	err := srv.DB.SelectContext(ctx, &newEntries, `
		SELECT
			fent.id,
			fent.source_id,
			fent.link,
			fent.published_on,
			fent.title,
			fent.description,
			fent.content
		FROM feed_entries AS fent
		LEFT JOIN feed_sent_emails AS fsem
			ON fsem.entry_id = fent.id AND fsem.to_email = ?
		WHERE fsem.id IS NULL
	`, notSentTo)
	if err != nil {
		return nil, err
	}

	return newEntries, nil
}
