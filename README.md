# Gator

Gator is a command-line RSS feed aggregator written in Go. Created as part of the [Buid a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang) course on [boot.dev](https://www.boot.dev).

## Requirements

- Go **1.26.4 or newer** (the version specified in `go.mod`)
- PostgreSQL, running locally or reachable from your machine
- [Goose](https://github.com/pressly/goose) to apply the database migrations

## Install the CLI

Install the latest version with Go:

```sh
go install github.com/novotnyradekcz/gator@latest
```

Make sure Go's bin directory is on your `PATH` (commonly `$(go env GOPATH)/bin`). The installed command is `gator`.

## Set up PostgreSQL and the schema

Create a PostgreSQL database for Gator. For example, with PostgreSQL's command-line tools:

```sh
createdb gator
```

From the root of a checkout of this repository, set the database URL and apply the migrations in order. The migrations in `sql/schema` create the tables Gator needs.

```sh
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "$DATABASE_URL" up
```

Change the username, password, host, port, and database in `DATABASE_URL` to match your PostgreSQL setup. If your password contains special characters, URL-encode them in the connection URL.

## Configure Gator

Gator reads its config from `~/.gatorconfig.json` (in your home directory). Create that file and set `db_url` to the same PostgreSQL connection URL:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

You can register a user after configuration; registration saves that user as the current user in the config file.

## Run Gator

Register a user, then add an RSS feed. Adding a feed also makes the current user follow it:

```sh
gator register alice
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
```

Start the feed aggregator in a terminal. It keeps running and checks feeds at the interval you provide; stop it with `Ctrl+C`:

```sh
gator agg 1m
```

In another terminal, browse the collected posts or inspect feeds and follows:

```sh
gator browse        # Show the 2 most recent posts by default
gator browse 10     # Show up to 10 posts
gator feeds         # List registered feeds
gator following     # List feeds followed by the current user
```

Other useful commands:

```sh
gator users                                     # List users
gator login alice                               # Switch the current user
gator follow "https://example.com/feed.xml"     # Follow an existing feed by URL
gator unfollow "https://example.com/feed.xml"   # Stop following a feed
```

`follow` requires that the feed has already been added with `addfeed`. Commands such as `addfeed`, `follow`, `unfollow`, `following`, and `browse` require a current user.
