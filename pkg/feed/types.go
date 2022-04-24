package feed

import "time"

type Feed struct {
	ID      int
	Link    string
	Title   string
	Authors string
}

type Entry struct {
	ID          int
	FeedID      int
	Link        string
	PublishedOn time.Time
	Title       string
	Description string
	Content     string
}
