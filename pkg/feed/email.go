package feed

import (
	"context"
	"fmt"
	"html"

	"github.com/azuline/presage/pkg/services"
)

func SendEntry(
	ctx context.Context,
	srv *services.Services,
	to string,
	entry EntryWithSourceTitle,
) error {
	if err := srv.Email.SendEmail(to, constructSubject(entry), constructBody(entry)); err != nil {
		return err
	}
	return recordSentEntry(ctx, srv, to, entry)
}

func constructSubject(entry EntryWithSourceTitle) string {
	return fmt.Sprintf("RSS: [%s] %s", entry.SourceTitle, entry.Title)
}

func constructBody(entry EntryWithSourceTitle) string {
	body := ""

	if entry.Title != "" {
		body += fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(entry.Title))
	}
	if entry.SourceTitle != "" {
		body += fmt.Sprintf("<h2>%s</h2>\n", html.EscapeString(entry.SourceTitle))
	}
	body += "Published on " + entry.PublishedOn + "<br><br>"

	body += fmt.Sprintf(
		"Link: <a href=\"%s\">%s</a><br><br><hr><br>",
		html.EscapeString(entry.Link),
		html.EscapeString(entry.Link),
	)

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

func BackfillSendingEntry(
	ctx context.Context,
	srv *services.Services,
	to string,
	entry EntryWithSourceTitle,
) error {
	return recordSentEntry(ctx, srv, to, entry)
}

func recordSentEntry(
	ctx context.Context,
	srv *services.Services,
	to string,
	entry EntryWithSourceTitle,
) error {
	_, err := srv.DB.ExecContext(
		ctx,
		"INSERT INTO feed_sent_emails (entry_id, to_email) VALUES (?, ?)",
		entry.ID, to,
	)
	if err != nil {
		return err
	}

	return nil
}
