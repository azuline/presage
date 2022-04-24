package feed

import (
	"context"
	"fmt"
	"html"

	"github.com/azuline/presage/pkg/services"
)

func SendEntry(ctx context.Context, srv *services.Services, to string, entry EntryWithSourceTitle) error {
	if err := srv.Email.SendEmail(to, constructSubject(entry), constructBody(entry)); err != nil {
		return err
	}

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

func constructSubject(entry EntryWithSourceTitle) string {
	return fmt.Sprintf("RSS: [%s] %s", entry.SourceTitle, entry.Title)
}

func constructBody(entry EntryWithSourceTitle) string {
	body := ""

	if entry.Title != "" {
		body += fmt.Sprintf("<h1>%s</h1><br>", html.EscapeString(entry.Title))
	}
	if entry.SourceTitle != "" {
		body += fmt.Sprintf("<h3>%s</h3><br><br>", html.EscapeString(entry.SourceTitle))
	}
	if !entry.PublishedOn.IsZero() {
		body += "Published on " + entry.PublishedOn.Format("02 Jan 06") + "<br><br>"
	}

	body += fmt.Sprintf("Link: <a href=\"%s\">%s</a><br><br>", html.EscapeString(entry.Link), html.EscapeString(entry.Link))

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
