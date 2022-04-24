CREATE TABLE feed_sources (
    id      INTEGER PRIMARY KEY,
    link    TEXT    UNIQUE NOT NULL,
    title   TEXT    NOT NULL,
    authors TEXT    NOT NULL
);

CREATE TABLE feed_entries (
    id           INTEGER  PRIMARY KEY,
    source_id    INTEGER  NOT NULL REFERENCES feed_sources(id),
    link         TEXT     UNIQUE NOT NULL,
    published_on DATETIME NOT NULL,
    title        TEXT     NOT NULL,
    description  TEXT     NOT NULL,
    content      TEXT     NOT NULL
);

CREATE TABLE feed_sent_emails (
    id       INTEGER  PRIMARY KEY,
    entry_id INTEGER  NOT NULL REFERENCES feed_entries(id),
    to_email TEXT     NOT NULL,
    sent_on  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
