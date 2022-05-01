.PHONY: build
build:
	earthly --use-inline-cache +filter-feed
	earthly --use-inline-cache +ui

gen-api:
	earthly --use-inline-cache  +proto 

clean:
	rm -rf bin/