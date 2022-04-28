all:
	earthly --use-inline-cache +filter-feed

proto: proto/*.proto
	earthly --use-inline-cache  +proto 
