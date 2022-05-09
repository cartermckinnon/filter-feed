# `filter-feed`

This project is:
1. A stateless, (aspirationally) transparent, filtering proxy for RSS, Atom, and JSON feeds.
2. A command-line tool for fetching and filtering feeds.
3. A user interface to manage a feed's filters.

I created this for podcasts. Your mileage may vary with other types of feeds -- that support comes purely from these awesome libraries that power this project:

- `github.com/mmcdole/gofeed`
- `github.com/gorilla/feeds`

## Usage

### On the web.

You can create a filtered feed URL from an existing feed using the public instance of `filter-feed` available at:

[`https://filter-feed.me`](https://filter-feed.me)

Please don't abuse this service 🙏.

### As a command-line tool.
```
filter-feed fetch URL FILTER (OVERRIDE)
```
Where:
- `URL` is the URL for the feed.
- `FILTER` is either:
  - A `FilterSpec`
  - A `FilterSpec`
  - A path to a JSON file beginning with `file://` (ex. `file://my-filter.json`), that contains a `FilterSpec` or `FilterSpecs`.
- (Optional) `OVERRIDE` is either:
  - A `OverrideSpec`
  - A `OverrideSpecs`
  - path to a JSON file, beginning with `file://` (ex. `file://my-override.json`), that contains a `OverrideSpec` or `OverrideSpecs`. 

### As a proxy.

Start the server:
```sh
filter-feed server --redis-address ":6379"
```

Or start the server in a container:
```sh
docker run --rm -p 8080:8080 ghcr.io/cartermckinnon/filter-feed server --redis-address "redis.host:6379"
```

Generate a URL on the `filter-feed.me` website, or with the command-line:
```sh
filter-feed url URL FILTER (OVERRIDE)
```

Where `URL`, `FILTER`, and `OVERRIDE` are as defined above.

The filtered URL's path will look something like:
```
/v1/ff/ChRodHRwczovL2NhdC5mZWVkL3JzcxIYCgVyZWdleBIILipjYXRzLioaBXRpdGxl
```

Just send a `GET` for that path to a `filter-feed` server, such as the one I run at `api.filter-feed.me`:
```
curl https://api.filter-feed.me/v1/ff/ChRodHRwczovL2NhdC5mZWVkL3JzcxIYCgVyZWdleBIILipjYXRzLioaBXRpdGxl
```

## Development

1. Install [Earthly](https://earthly.dev).

2. Compile the protocol buffers:
```sh
make proto
```

3. Build locally:
```sh
go build
```

4. Build for release:
```sh
make
```
