package override

import (
	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

type Override interface {
	Apply(*gofeed.Feed) bool
	GetSpec() *api.OverrideSpec
}
