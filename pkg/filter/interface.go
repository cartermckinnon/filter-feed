package filter

import (
	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type Filter interface {
	Matches(*gofeed.Item) bool
	GetSpec() *api.FilterSpec
}
