# Filter feed

A stateless, filtering proxy for RSS, Atom, and JSON feeds.

## Development

0. Install [Earthly](https://earthly.dev).

1. Compile the protocol buffers:
```sh
earthly +proto
```

2. Build:
```sh
earthly +filter-feed
```

## Usage

```sh
docker run --rm -p 8080:8080 ghcr.io/cartermckinnon/filter-feed
```