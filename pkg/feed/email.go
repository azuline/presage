package feed

import (
	"context"
	"fmt"

	"github.com/azuline/presage/pkg/services"
)

func SendEntry(ctx context.Context, srv *services.Services, to string, entry EntryWithSourceTitle) error {
	if err := srv.Email.SendEmail(to, constructSubject(entry), constructBody(entry)); err != nil {
		return err
	}

	_, err := srv.DB.ExecContext(
		ctx,
		"INSERT INTO feeds_sent_emails (entry_id, sent_to) VALUES (?, ?)",
		entry.ID, to,
	)
	if err != nil {
		return err
	}

	return nil
}

func constructSubject(entry EntryWithSourceTitle) string {
	return fmt.Sprintf("RSS: [%s] %s", entry.SourceTitle, entry.Title)
}

func constructBody(entry EntryWithSourceTitle) string {
	body := ""
	if !entry.PublishedOn.IsZero() {
		body += "Published on " + entry.PublishedOn.Format("02 Jan 06") + "\r\n\r\n"
	}

	body += fmt.Sprintf("<a href=\"%s\">%s</a>\r\n\r\n", entry.Link, entry.Link)

	switch {
	case entry.Description != "":
		body += entry.Description
	case entry.Content != "":
		body += entry.Content
	default:
		body += "No preview."
	}

	return body
}
