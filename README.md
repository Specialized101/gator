# Gator - A RSS Aggregator in GO

## How to use

1. You will need GO and Postgres installed (linux only, see [webi](https://webinstall.dev/) for MAC/Windows)
```
curl -sS https://webi.sh/golang | sh; \
curl -sS https://webi.sh/postgres | sh; \
source ~/.config/envman/PATH.env
```

2. Clone this repo
```
git clone https://github.com/Specialized101/gator
cd gator
```

3. Install the CLI tool
```
go install
```

4. Create the configuration file to hold the database connection string
```
echo '{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}' > ~/.gatorconfig.json
```

5. Setup the postgres user password to 'postgres'
```
sudo passwd postgres
```

6. Create the gator database
```
sudo -u postgres psql
CREATE DATABASE gator;
\c gator
ALTER USER postgres PASSWORD 'postgres';
```

7. Install goose
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

8. Run the database migrations
```
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

9. Run the program
```
gator <command> [arguments]
```

## List of available commands
- login: gator login <name>       -> login to gator as name
- register: gator register <name> -> register name in the gator database
- users: gator users              -> list all registered users
- addfeed: gator addfeed <name> <url> -> add a rss feed
- agg: gator agg <duration> -> fetch rss feeds every duration (examples of duration: 1s, 2m, 1h, etc)
    example: gator agg 5s will fetch from added rss feeds every 5 seconds
- feeds: gator feeds -> list all feeds from the database
- follow: gator follow <url> -> subscribe to a rss feed
- unfollow: gator unfollow <url> -> unsubcribe from rss feed
- following: gator following -> list feeds you are currently subscribed to
- browse: gator browse <limit> -> show limit number of posts from the followed feeds