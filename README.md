# `filter-feed`

This project produces one binary that contains:
1. A stateless, (aspirationally) transparent, filtering proxy for RSS, Atom, and JSON feeds.
2. A command-line tool for fetching and filtering feeds.

## Usage

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

Generate a URL:
```sh
filter-feed url URL FILTER (OVERRIDE)
```

Where `URL`, `FILTER`, and `OVERRIDE` are as defined above.

It will look something like:
```
/v1/f/ChRodHRwczovL2NhdC5mZWVkL3JzcxIYCgVyZWdleBIILipjYXRzLioaBXRpdGxl
```

Just send a `GET` for that path to a `filter-feed` server, such as the one I run at `api.filter-feed.me`:
```
curl https://api.filter-feed.me/v1/f/ChRodHRwczovL2NhdC5mZWVkL3JzcxIYCgVyZWdleBIILipjYXRzLioaBXRpdGxl
```

For example, you could add such a link to your favorite podcast player.

🔨👷‍♂️🚧 I intend to provide a form at `https://filter-feed.me/` that can generate feed URLs; but this isn't finished yet.

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
