package override

import (
	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/mmcdole/gofeed"
)

func ApplyOverrides(feed *gofeed.Feed, specs []*api.OverrideSpec) error {
	modified := false
	for _, spec := range specs {
		o, err := buildOverride(spec)
		if err != nil {
			return err
		}
		modified = modified || o.Apply(feed)
	}
	return nil
}
