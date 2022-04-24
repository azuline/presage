# presage

`presage` is a tool that scrapes RSS feeds and sends out an email for new
articles.

Scraped articles are stored in a SQLite database. All outbound emails are
recorded in the SQLite database to ensure only new articles are sent.

Outbound emails are recorded as a `(sent_to_email, article_link)` tuple.
"New" emails are tracked per email address. See [Backfill](#backfill) to
backfill outbound email records for a new email address.

**You probably want to backfill first when downloading.**

## Installation

Due to this project being personal, we don't distribute any binaries right now.
The binary can be compiled natively with the standard Go 1.18 toolchain and the
`make build` Makefile rule.

The Go toolchain is available with the provided `shell.nix`, or can be
installed separately.

## Usage

This tool is intended to be run periodically via a cronjob.

This tool must be called with all the following parameters and environment
variables. We send outbound email via SMTP and store state in a SQLite database
file.

```
export SMTP_USER=username@ema.il
export SMTP_PASS=hopeyouknowit
export SMTP_HOST=lol.find.your.own
export SMTP_PORT=587
export DATABASE_PATH=./presage.sqlite

presage \
  -send-to recipient@ema.il \
  -feeds-list /path/to/feeds.txt
```

### Environment Variable File

We support reading environment variables from a file. This defaults to `.env`,
but can be customized by the `-env-file` parameter.

### Feeds List Format

The file referenced by `-feeds-list` must be formatted like so:

```
http://feed.one/something/rss.xml
http://feed.two/again/atom.xml
http://feed.three/whatever/hopeitsxml
```

### Backfill

A `-backfill` flag takes all articles and records an outbound
`(sent_to_email, article_link)` for each article. However, it does not send an
email. This can be used to avoid sending every prior article to a new email
address.

### Dry Run

The optional `-dry-run` flag does everything up until sending and logging the
emails, at which point it stops and does not continue. No records of outbound
emails are recorded, nor are any sent.

## Development

There is an included `shell.nix` file for setting up the tooling needed to work
on this project. We recommend using `direnv` to auto-enter the environment.

## License

```
presage :: a rss feed scraper and email tool

Copyright (C) 2022 blissful <blissful@sunsetglow.net>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free
Software Foundation, either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.  See the GNU Affero General Public License for more
details.

You should have received a copy of the GNU Affero General Public License along
with this program.  If not, see <https://www.gnu.org/licenses/>.
```
