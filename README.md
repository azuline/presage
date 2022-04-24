# presage

`presage` is a tool that scrapes RSS feeds and sends out an email for new
articles via Sendgrid.

Sendgrid's free plan supports up to 100 emails a day. I hope none of us have
more than 100 new RSS articles a day!

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
