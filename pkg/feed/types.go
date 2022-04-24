package feed

import "time"

type Feed struct {
	ID      int    `db:"id"`
	Link    string `db:"link"`
	Title   string `db:"title"`
	Authors string `db:"authors"`
}

type Entry struct {
	ID          int       `db:"id"`
	FeedID      int       `db:"source_id"`
	Link        string    `db:"link"`
	PublishedOn time.Time `db:"published_on"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Content     string    `db:"content"`
}

type SentEmail struct {
	ID      int       `db:"id"`
	EntryID int       `db:"entry_id"`
	ToEmail string    `db:"to_email"`
	SentOn  time.Time `db:"sent_on"`
}
