package filter

import (
	api "github.com/cmckn/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type Filter interface {
	Matches(*gofeed.Item) bool
	GetSpec() *api.FilterSpec
}
