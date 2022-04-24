package feed

import (
	"context"

	"github.com/pkg/errors"

	"github.com/azuline/presage/pkg/services"
)

type EntryWithSourceTitle struct {
	Entry
	SourceTitle string `db:"source_title"`
}

// ReadNewEntries returns all feed entries that have not yet been sent out over
// email.
func ReadNewEntries(ctx context.Context, srv *services.Services, notSentTo string) ([]EntryWithSourceTitle, error) {
	var newEntries []EntryWithSourceTitle

	err := srv.DB.SelectContext(ctx, &newEntries, `
		SELECT
			fent.id,
			fent.source_id,
			fent.link,
			fent.published_on,
			fent.title,
			fent.description,
			fent.content,
			fsrc.title AS source_title
		FROM feed_entries AS fent
		JOIN feed_sources AS fsrc ON fsrc.id = fent.source_id
		WHERE NOT EXISTS (
			SELECT 1 FROM feed_sent_emails AS fsem
			WHERE fsem.entry_id = fent.id AND fsem.to_email = ?
		)
	`, notSentTo)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading new entries")
	}

	return newEntries, nil
}
